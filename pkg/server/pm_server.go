package server

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	TypeUnixSocket = "unix"
	TypeTcp        = "tcp"

	EnvPreforkChildKey = "ZINC_PREFORK_CHILD"
	EnvExecTimeout     = "ZINC_EXEC_TIMEOUT"
	EnvNamedPipeID     = "ZINC_PIPE_ID"
	EnvPreforkChildVal = "OK"

	NamedPipe = "named-pipe"

	WORKER_STATE_IDLE    uint8 = 1 << 0
	WORKER_STATE_BUSY    uint8 = 1 << 1
	WORKER_STATE_STOPPED uint8 = 1 << 2

	defaultTimeout = 60 // default is 60s
)

// ZnPMServer - use "1 request, 1 process" mode like what PHP-FPM is doing
// To start a ZnPMServer, we must call StartMaster() first to create network connection, then spawn another process to call StartWorker() to run actual `reqHandler` logic.
type ZnPMServer struct {
	childProcManager
	reqHandler http.HandlerFunc
}

type childProcManager struct {
	//// internal properites for child proc management
	childs map[int]workerState

	refCount   int
	addChan    chan workerState
	updateChan chan workerState
	delChan    chan int
}

type workerState struct {
	pid   int
	cmd   *exec.Cmd
	state uint8
}

type ZnPMServerConfig struct {
	InitProcs int
	// maximum procs the worker could create
	MaxProcs int
	Timeout  int
}

func NewZnPMServer(reqHandler http.HandlerFunc) *ZnPMServer {
	return &ZnPMServer{
		childProcManager: childProcManager{
			childs:     make(map[int]workerState),
			addChan:    make(chan workerState),
			updateChan: make(chan workerState),
			delChan:    make(chan int),
			refCount:   0,
		},
		reqHandler: reqHandler,
	}
}

func (zns *ZnPMServer) Start(connUrl string, cfg ZnPMServerConfig) error {
	childWorkerFlag := false

	// check if --child-worker flag exists
	for _, arg := range os.Args {
		if arg == "--child-worker" {
			childWorkerFlag = true
			break
		}
	}

	///// run child worker if  --child-worker = true & preForkChild env is "OK"
	if childWorkerFlag && os.Getenv(EnvPreforkChildKey) == EnvPreforkChildVal {
		// start child worker to handle requests
		if err := zns.StartWorker(); err != nil {
			return err
		}
	} else {
		//// otherwise, just listen to the server
		if err := zns.StartMaster(connUrl, cfg); err != nil {
			return err
		}
	}

	return nil
}

////////////////////////////
/////// MASTER logic ///////
////////////////////////////
func (zns *ZnPMServer) StartMaster(connUrl string, cfg ZnPMServerConfig) error {
	network, address, err := parseConnUrl(connUrl)
	if err != nil {
		return err
	}

	log.Print("即将监听URL：", address)
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	// since ZnFPMServer only accepts tcp and unix, the net.Listener MUST
	// be TCPListener
	ln := l.(*net.TCPListener)

	log.Print("即将打开父-子进程通信通道")
	p, err := CreateNamedPipe()
	if err != nil {
		log.Fatalf("CreatePipe: failed - %s", err)
		return err
	}

	// kill child processes when master exits
	defer func() {
		var s sync.RWMutex
		s.Lock()
		for _, procState := range zns.childs {
			proc := procState.cmd
			if err := proc.Process.Kill(); err != nil {
				if !errors.Is(err, os.ErrProcessDone) {
					log.Fatalf("spawnProcs: failed to kill child: %v", err)
				}
			}
		}
		s.Unlock()
	}()

	//// read named pipe data to recv msg from child process
	go zns.readNamedPipe(p)

	//// maintain child state (DO NOT UPDATE child data directly!)
	go zns.maintainChildState(cfg, ln, p)

	for i := 0; i < cfg.InitProcs; i++ {
		if err := zns.spawnProcess(cfg, ln, p); err != nil {
			return err
		}
	}

	// register signal handler
	// Unix sockets must be unlink()ed before being reused again.

	// Handle common process-killing signals so we can gracefully shut down:
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	// Wait for a SIGINT or SIGKILL:
	<-sigc
	log.Print("正在关闭服务器...")
	// Stop listening (and unlink the socket if unix type):
	return l.Close()
}

//// fork child processes
func (zns *ZnPMServer) spawnProcess(cfg ZnPMServerConfig, l *net.TCPListener, p *pipe) error {
	// prepare net.Conn file to transfer to child processes
	lf, err := l.File()
	if err != nil {
		return err
	}
	// close fd only effective on current process only
	syscall.CloseOnExec(int(lf.Fd()))

	// spawn new child processes
	cmd := exec.Command(os.Args[0], "--child-worker")

	// add prefork child flag into child proc env
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("%s=%s", EnvPreforkChildKey, EnvPreforkChildVal),
		fmt.Sprintf("%s=%d", EnvExecTimeout, cfg.Timeout),
		fmt.Sprintf("%s=%s", EnvNamedPipeID, GetPipeID(p)),
	)

	// pass connection FD to child process as ExtraFile
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.ExtraFiles = []*os.File{lf}

	if err := cmd.Start(); err != nil {
		return err
	}
	// store child process
	pid := cmd.Process.Pid

	// register new child process to workerState
	zns.addChan <- workerState{
		pid:   pid,
		state: WORKER_STATE_IDLE,
		cmd:   cmd,
	}
	// send msg to channel when cmd ends running
	go func() {
		cmd.Wait()
		// after cmd is done, send pid to del channel
		zns.delChan <- pid
	}()

	return nil
}

func (zns *ZnPMServer) readNamedPipe(pipe *pipe) {
	pipeReader, err := OpenNamedPipeReader(pipe)
	if err != nil {
		log.Fatal("[PARENT] Open named pipe file error:", err)
		return
	}

	var buf = make([]byte, 5)
	for {
		var state uint8
		var pid int
		// read packet
		if err := ReadDataFromNamedPipe(pipeReader, buf); err != nil {
			log.Fatalf("[PARENT] read buffer failed: %s", err)
			continue
		}

		// digest state & pid
		pid = int(binary.BigEndian.Uint32(buf))
		state = buf[4]

		zns.updateChan <- workerState{
			pid:   pid,
			state: state,
			cmd:   nil,
		}
	}
}

// summon all writing actions into one goroutine to ensure thread-safe on writing.
func (zns *ZnPMServer) maintainChildState(cfg ZnPMServerConfig, ln *net.TCPListener, p *pipe) {
	for {
		select {
		case aw := <-zns.addChan:
			zns.childs[aw.pid] = aw
			zns.refCount = len(zns.childs)
		case uw := <-zns.updateChan:
			if oldState, ok := zns.childs[uw.pid]; ok {
				zns.childs[uw.pid] = workerState{
					pid:   uw.pid,
					state: uw.state,
					cmd:   oldState.cmd,
				}
			}
			// spawn more process (total procs not exceed the number of `maxProcs`)
			// when there's no idle process
			// count the number of idle processes
			hasIdleProcs := false
			for _, w := range zns.childs {
				if w.state == WORKER_STATE_IDLE {
					hasIdleProcs = true
					break
				}
			}

			// all procs are busy (or stopped), spawn process first
			if !hasIdleProcs {
				currentNum := zns.refCount
				finalProcNum := currentNum + 10

				if finalProcNum > cfg.MaxProcs {
					finalProcNum = cfg.MaxProcs
				}

				addNum := finalProcNum - currentNum
				zns.refCount = finalProcNum
				go func() {
					for i := 0; i < addNum; i++ {
						if err := zns.spawnProcess(cfg, ln, p); err != nil {
							log.Fatalf("启动子进程失败：%v", err)
							break
						}
					}
				}()
			}
		case pid := <-zns.delChan:
			delete(zns.childs, pid)
			zns.refCount -= 1
			// check if the number of current existing procs is lower than `initProcs`
			if zns.refCount < cfg.InitProcs {
				numsToSpawn := cfg.InitProcs - zns.refCount
				zns.refCount += numsToSpawn
				// spawn more processes to ensure minimum proc number
				go func() {
					for i := 0; i < numsToSpawn; i++ {
						// add time delay to avoid non-stop reboot too quick
						time.Sleep(100 * time.Millisecond)
						if err := zns.spawnProcess(cfg, ln, p); err != nil {
							log.Fatalf("启动子进程失败：%v", err)
							break
						}
					}
				}()
			}
		}
	}
}

////////////////////////////
/////// WORKER logic ///////
////////////////////////////

//// USE ANOTHER PROCESS TO START!
func (zns *ZnPMServer) StartWorker() error {
	pid := os.Getpid()
	lf := os.NewFile(uintptr(3), fmt.Sprintf("listener-%d", pid))
	defer lf.Close()

	// read pipeID from env
	p := NewPipe(os.Getenv(EnvNamedPipeID))
	pipeWriter, err := OpenNamedPipeWriter(p)
	if err != nil {
		log.Fatalf("[CHILD] OpenNamedPipeWriter failed: %s", err.Error())
		return err
	}
	defer pipeWriter.Close()

	lc, err := net.FileListener(lf)
	if err != nil {
		log.Fatalf("[CHILD] Inherit Parent TCPListener failed: %s", err.Error())
		return err
	}
	// get execution timeout
	timeout := defaultTimeout
	if t, err := strconv.Atoi(os.Getenv(EnvExecTimeout)); err == nil {
		timeout = t
	}

	// kill child process when parent process exits
	go func() {
		// if it is equal to 1 (init process ID),
		// it indicates that the master process has exited
		const watchInterval = 500 * time.Millisecond
		for range time.NewTicker(watchInterval).C {
			if os.Getppid() == 1 {
				os.Exit(1) //nolint:revive // Calling os.Exit is fine here in the prefork
			}
		}
	}()

	for {
		conn, err := lc.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %v", err)
		}

		// set busy state
		zns.writeProcState(pipeWriter, WORKER_STATE_BUSY)

		// Wrap the connection in a bufio.Reader to read the HTTP request
		bufReader := bufio.NewReader(conn)
		req, err := http.ReadRequest(bufReader)
		if err != nil {
			return fmt.Errorf("read HTTP request error: %v", err)
		}

		waitSig := make(chan int)

		// handle request
		go func() {
			w := NewRespWriter(conn)
			zns.reqHandler(w, req)

			waitSig <- 1
		}()

		select {
		case <-waitSig:
			zns.writeProcState(pipeWriter, WORKER_STATE_IDLE)
		case <-time.After(time.Duration(timeout) * time.Second):
			zns.writeProcState(pipeWriter, WORKER_STATE_STOPPED)
			// wait for a while before exiting the process
			time.Sleep(100 * time.Millisecond)
			// exit the process directly
			os.Exit(1)
		}

		conn.Close()
	}
}

func (zns *ZnPMServer) writeProcState(writer *os.File, state uint8) {
	pid := os.Getpid()
	buf := make([]byte, 5)
	// write state msg:
	// +-------------------------------------------+
	// | pid[3] | pid[2] | pid[1] | pid[0] | state |
	// +-------------------------------------------+

	binary.BigEndian.PutUint32(buf, uint32(pid))
	buf[4] = state

	// write state msg to notify parent process
	writer.Write(buf)
}

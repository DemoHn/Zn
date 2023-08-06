package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	TypeUnixSocket     = "unix"
	TypeTcp            = "tcp"
	EnvZincAdapter     = "ZINC_ADAPTER"
	ZincAP_Playground  = "playground"
	ZincAP_HTTPHandler = "http_handler"
	EnvPreforkChildKey = "ZINC_PREFORK_CHILD"
	EnvExecTimeout     = "ZINC_EXEC_TIMEOUT"
	EnvPreforkChildVal = "OK"

	NamedPipe = "named-pipe"
)

type ZnServer struct {
	Network string
	// used fo network = "tcp" or "unix"
	Address string

	Config ZnServerConfig

	//// internal properites for child proc management
	childs map[int]workerState

	refCount   int
	addChan    chan workerState
	updateChan chan workerState
	delChan    chan int
}

type ZnServerConfig struct {
	InitProcs int
	// maximum procs the worker could create
	MaxProcs int
	Timeout  int
}

type workerState struct {
	pid   int
	cmd   *exec.Cmd
	state uint8
}

// e.g.: unix:///home/demohn/test.sock
func NewFromURL(connUrl string, config ZnServerConfig) (*ZnServer, error) {
	u, err := url.Parse(connUrl)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case TypeUnixSocket:
		return newZnServer(TypeUnixSocket, u.Path, config), nil
	case TypeTcp:
		return newZnServer(TypeTcp, u.Host, config), nil
	}
	return nil, fmt.Errorf("不支持的协议：%s", u.Scheme)
}

func newZnServer(network string, address string, config ZnServerConfig) *ZnServer {
	return &ZnServer{
		Network:    network,
		Address:    address,
		Config:     config,
		childs:     make(map[int]workerState),
		addChan:    make(chan workerState),
		updateChan: make(chan workerState),
		delChan:    make(chan int),
		refCount:   0,
	}
}

func (zns *ZnServer) StartMaster() error {
	cfg := zns.Config

	log.Print("即将监听URL：", zns.Address)
	l, err := net.Listen(zns.Network, zns.Address)
	if err != nil {
		return err
	}

	// since ZnServer only accepts tcp and unix, the net.Listener MUST
	// be TCPListener
	ln := l.(*net.TCPListener)

	log.Print("即将打开父-子进程通信通道")
	if err := syscall.Mkfifo(NamedPipe, 0666); err != nil {
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
	go zns.readNamedPipe(NamedPipe)

	// open named pipe write to child process
	pipeWriter, err := os.OpenFile(NamedPipe, os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	//// maintain child state (DO NOT UPDATE child data directly!)
	go zns.maintainChildState(ln, pipeWriter)

	for i := 0; i < cfg.InitProcs; i++ {
		if err := zns.spawnProcess(ln, pipeWriter); err != nil {
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
func (zns *ZnServer) spawnProcess(l *net.TCPListener, pipeWriter *os.File) error {
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
		fmt.Sprintf("%s=%d", EnvExecTimeout, zns.Config.Timeout),
	)

	// pass connection FD to child process as ExtraFile
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.ExtraFiles = []*os.File{lf, pipeWriter}

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

func (zns *ZnServer) readNamedPipe(pipeFile string) {
	pipeReader, err := os.OpenFile(pipeFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("[PARENT] Open named pipe file error:", err)
		return
	}
	os.Remove(pipeFile)

	var buf = make([]byte, 5)
	for {
		var state uint8
		var pid int
		// read packet
		_, err := pipeReader.Read(buf)
		if err != nil {
			log.Fatal("[PARENT] read buffer failed")
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
func (zns *ZnServer) maintainChildState(ln *net.TCPListener, pipeWriter *os.File) {
	cfg := zns.Config
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

				if finalProcNum > zns.Config.MaxProcs {
					finalProcNum = zns.Config.MaxProcs
				}

				addNum := finalProcNum - currentNum
				zns.refCount = finalProcNum
				go func() {
					fmt.Println("update spawn: ", addNum)
					for i := 0; i < addNum; i++ {
						if err := zns.spawnProcess(ln, pipeWriter); err != nil {
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
					fmt.Println("delete spawn: ", numsToSpawn)

					for i := 0; i < numsToSpawn; i++ {
						// add time delay to avoid non-stop reboot too quick
						time.Sleep(100 * time.Millisecond)
						if err := zns.spawnProcess(ln, pipeWriter); err != nil {
							log.Fatalf("启动子进程失败：%v", err)
							break
						}
					}
				}()
			}
		}
	}
}

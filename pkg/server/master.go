package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

const (
	TypeUnixSocket     = "unix"
	TypeTcp            = "tcp"
	EnvZincAdapter     = "ZINC_ADAPTER"
	ZincAP_Playground  = "playground"
	ZincAP_HTTPHandler = "http_handler"
	EnvPreforkChildKey = "ZINC_PREFORK_CHILD"
	EnvPreforkChildVal = "OK"
)

type ZnServer struct {
	Network string
	// used fo network = "tcp" or "unix"
	Address string
}

type ZnServerConfig struct {
	InitProcs int
	// maximum procs the worker could create
	MaxProcs int
}

type worker struct {
	pid int
	err error
}

// e.g.: unix:///home/demohn/test.sock
func NewFromURL(connUrl string) (*ZnServer, error) {
	u, err := url.Parse(connUrl)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case TypeUnixSocket:
		return &ZnServer{
			Network: TypeUnixSocket,
			Address: u.Path,
		}, nil
	case TypeTcp:
		return &ZnServer{
			Network: TypeTcp,
			Address: u.Host,
		}, nil
	}
	return nil, fmt.Errorf("不支持的协议：%s", u.Scheme)
}

func (zns *ZnServer) StartMaster(cfg ZnServerConfig) error {
	var l net.Listener
	var err error

	log.Print("即将监听URL：", zns.Address)
	l, err = net.Listen(zns.Network, zns.Address)
	if err != nil {
		return err
	}

	log.Print("即将打开父-子进程通信通道")
	namedPipe, err := createNamedPipe()
	if err != nil {
		return err
	}

	childs := make(map[int]*exec.Cmd)
	channel := make(chan worker, cfg.MaxProcs)

	// kill child processes when master exits
	defer func() {
		for _, proc := range childs {
			if err := proc.Process.Kill(); err != nil {
				if !errors.Is(err, os.ErrProcessDone) {
					fmt.Printf("prefork: failed to kill child: %v", err)
				}
			}
		}
	}()

	//// read named pipe data to recv msg from child process
	go readNamedPipe(namedPipe)

	pipeWriter, err := os.OpenFile(namedPipe, os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	// Since ZnServer only accepts tcp and unix, the net.Listener MUST
	// be TCPListener
	for i := 0; i < cfg.InitProcs; i++ {
		cmd, err := prefork(l.(*net.TCPListener), channel, pipeWriter)
		if err != nil {
			return err
		}
		pid := cmd.Process.Pid
		childs[pid] = cmd
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

func createNamedPipe() (string, error) {
	// generate random named pipe name
	pipeFile := fmt.Sprintf("pipe%d", rand.Intn(100000))

	// make name pipe
	if err := syscall.Mkfifo(pipeFile, 0666); err != nil {
		return "", err
	}

	return pipeFile, nil
}

//// fork child processes
func prefork(l *net.TCPListener, waitMsg chan worker, pipeWriter *os.File) (*exec.Cmd, error) {
	// prepare net.Conn file to transfer to child processes
	lf, err := l.File()
	if err != nil {
		return nil, err
	}
	// close fd only effective on current process only
	syscall.CloseOnExec(int(lf.Fd()))

	// spawn new child processes
	cmd := exec.Command(os.Args[0], "--child-worker")

	// add prefork child flag into child proc env
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("%s=%s", EnvPreforkChildKey, EnvPreforkChildVal),
	)

	// pass connection FD to child process as ExtraFile
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.ExtraFiles = []*os.File{lf, pipeWriter}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	// store child process
	pid := cmd.Process.Pid

	// send msg to channel when
	go func() {
		waitMsg <- worker{pid, cmd.Wait()}
	}()

	return cmd, nil
}

func readNamedPipe(pipeFile string) {
	pipeReader, err := os.OpenFile(pipeFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
		return
	}
	os.Remove(pipeFile)

	var pid int32
	for {
		// read packet
		binary.Read(pipeReader, binary.BigEndian, &pid)
		fmt.Println("read info", pid)
	}
}

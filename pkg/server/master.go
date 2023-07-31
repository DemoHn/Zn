package server

import (
	"fmt"
	"log"
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

func (zns *ZnServer) Listen() error {
	var l net.Listener
	var err error

	log.Print("即将监听URL：", zns.Address)
	l, err = net.Listen(zns.Network, zns.Address)
	if err != nil {
		return err
	}
	// Since ZnServer only accepts tcp and unix, the net.Listener MUST
	// be TCPListener
	if err := prefork(l.(*net.TCPListener), 3); err != nil {
		return err
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
	l.Close()
	return nil
	// os.Exit(0)

	//	log.Print("开始监听服务...")
	//	return fcgi.Serve(l, fcgiHandler)
}

//// fork child processes
func prefork(l *net.TCPListener, n int) error {
	// prepare net.Conn file to transfer to child processes
	lf, err := l.File()
	if err != nil {
		return err
	}
	// close fd only effective on current process only
	syscall.CloseOnExec(int(lf.Fd()))

	// spawn new child processes
	for i := 0; i < n; i++ {
		cmd := exec.Command(os.Args[0], "--child-worker")

		// add prefork child flag into child proc env
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("%s=%s", EnvPreforkChildKey, EnvPreforkChildVal),
		)

		// pass connection FD to child process as ExtraFile
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.ExtraFiles = []*os.File{lf}

		if err := cmd.Start(); err != nil {
			return err
		}
	}

	return nil
}

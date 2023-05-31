package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

const (
	TypeUnixSocket = "unix"
	TypeTcp        = "tcp"
)

type ZnServer struct {
	Network string
	// used fo network = "tcp" or "unix"
	Address string
}

type FCGIHandler struct{}

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

	// register signal handler
	// Unix sockets must be unlink()ed before being reused again.

	// Handle common process-killing signals so we can gracefully shut down:
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal) {
		// Wait for a SIGINT or SIGKILL:
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		// Stop listening (and unlink the socket if unix type):
		l.Close()
		// And we're done:
		os.Exit(0)
	}(sigc)

	log.Print("开始监听服务...")
	return fcgi.Serve(l, &FCGIHandler{})
}

func (f *FCGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<h2>This is HTTP Text</h2> body: <h3>%s</h3>", string("Hello World"))))
}

func listenToUnixSocket(socketPath string) (net.Listener, error) {
	// Check if the socket file exists
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		// Create a new socket file
		return net.Listen("unix", socketPath)
	}

	// Check if the socket is occupied by another process
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok {
			if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
				if sysErr.Err == syscall.ECONNREFUSED {
					return net.Listen("unix", socketPath)
				} else {
					return nil, fmt.Errorf("failed to connect to socket: %v", err)
				}
			} else {
				return nil, fmt.Errorf("failed to connect to socket: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to connect to socket: %v", err)
		}
	} else {
		conn.Close()
		return nil, fmt.Errorf("socket file is already occupied")
	}
}

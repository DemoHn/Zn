package server

import (
	"net"
	"net/http"
	"net/http/fcgi"
)

const (
	TypeUnixSocket = "unix"
	TypeTcp        = "tcp"
)

type ZnServer struct {
	Network string
	// used for network = "unix"
	SockFile string
	// used fo network = "tcp"
	Address string
}

type FCGIHandler struct{}

// e.g.: unix:///home/demohn/test.sock
func NewServer(url string) *ZnServer {
	zns := ZnServer{}

	return &zns
}

func (zns *ZnServer) Listen() error {
	addr := zns.SockFile
	if zns.Network == TypeTcp {
		addr = zns.Address
	}
	//// TODO: rewrite unix sock correctly

	l, err := net.Listen(zns.Network, addr)
	if err != nil {
		return err
	}
	return fcgi.Serve(l, &FCGIHandler{})
}

func (f *FCGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TBD...
}

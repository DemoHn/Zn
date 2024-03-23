package server

import (
	"log"
	"net"
	"net/http"
)

// NOTE: although we name it 'thread' server, we actually use goroutine
// as low level thread model instead of traditional "threads"
type ZnThreadServer struct {
	reqHandler http.HandlerFunc
}

func (zns *ZnThreadServer) Start(connUrl string) error {
	network, address, err := parseConnUrl(connUrl)
	if err != nil {
		return err
	}

	log.Print("即将监听URL：", address)
	ln, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	return http.Serve(ln, zns.reqHandler)
}

func NewZnThreadServer(reqHandler http.HandlerFunc) *ZnThreadServer {
	return &ZnThreadServer{
		reqHandler: reqHandler,
	}
}

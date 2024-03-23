package server

import (
	"log"
	"net"
	"net/http"
)

type ZnHTTPServer struct {
	reqHandler http.HandlerFunc
}

func (zns *ZnHTTPServer) Start(connUrl string) error {
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

func NewZnHTTPServer(reqHandler http.HandlerFunc) *ZnHTTPServer {
	return &ZnHTTPServer{
		reqHandler: reqHandler,
	}
}

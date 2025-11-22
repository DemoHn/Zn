package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// NOTE: although we name it 'thread' server, we actually use goroutine
// as low level thread model instead of traditional "threads"
type ZnThreadServer struct {
	reqHandler http.Handler
}

func (zns *ZnThreadServer) SetHandler(handler http.Handler) {
	zns.reqHandler = handler
}

func (zns *ZnThreadServer) Start(connUrl string) error {
	if zns.reqHandler == nil {
		return fmt.Errorf("处理逻辑未设置，需使用 setHandler() 配置")
	}

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

func NewZnThreadServer() *ZnThreadServer {
	return &ZnThreadServer{}
}

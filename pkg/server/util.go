package server

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

type RespWriter struct {
	writer *bufio.Writer
	header http.Header
}

func NewRespWriter(conn net.Conn) *RespWriter {
	return &RespWriter{
		writer: bufio.NewWriter(conn),
		header: http.Header{},
	}
}

//// fit http.ResponseWriter interface
func (w *RespWriter) Header() http.Header {
	return w.header
}

func (w *RespWriter) WriteHeader(statusCode int) {
	headString := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))

	w.writer.WriteString(headString)
	// write header (e.g. Content-Type: text/plain; charset="utf-8")
	w.header.Write(w.writer)
	// write \r\n
	w.writer.WriteString("\r\n")
}

func (w *RespWriter) Write(data []byte) (int, error) {
	nn, err := w.writer.Write(data)
	if err != nil {
		return nn, err
	}

	if err := w.writer.Flush(); err != nil {
		return nn, err
	}

	return nn, nil
}

//// parseConnUrl -> network + address
// currently support: tcp://, unix://
func parseConnUrl(connUrl string) (string, string, error) {
	u, err := url.Parse(connUrl)
	if err != nil {
		return "", "", err
	}
	switch u.Scheme {
	case TypeUnixSocket:
		return TypeUnixSocket, u.Path, nil
	case TypeTcp:
		return TypeTcp, u.Host, nil
	default:
		return "", "", fmt.Errorf("不支持的协议：%s", u.Scheme)
	}
}

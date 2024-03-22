package server

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
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

func (w *RespWriter) WriteSuccessData(body string) error {
	// declare content-type = text/plain (not json/file/...)
	w.addHeader("Content-Type", "text/plain; charset=\"utf-8\"")
	// http status: 200 OK
	w.writeHeader(http.StatusOK)
	// write resp body
	w.writer.WriteString(body)

	return w.writer.Flush()
}

func (w *RespWriter) WriteErrorData(reason error) error {
	// declare content-type = text/plain (not json/file/...)
	w.addHeader("Content-Type", "text/plain; charset=\"utf-8\"")
	// http status: 500 Internal Error
	w.writeHeader(http.StatusInternalServerError)
	// write resp body
	w.writer.WriteString(reason.Error())

	return w.writer.Flush()
}

func (w *RespWriter) addHeader(key string, value string) {
	w.header.Add(key, value)
}

// write HTTP/1.1 header data - including status code, header key/value
func (w *RespWriter) writeHeader(statusCode int) {
	headString := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))

	w.writer.WriteString(headString)
	// write header (e.g. Content-Type: text/plain; charset="utf-8")
	w.header.Write(w.writer)
	// write \r\n
	w.writer.WriteString("\r\n")
}

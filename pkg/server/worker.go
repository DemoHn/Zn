package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	WORKER_STATE_IDLE    uint8 = 1 << 0
	WORKER_STATE_BUSY    uint8 = 1 << 1
	WORKER_STATE_STOPPED uint8 = 1 << 2

	defaultTimeout = 60 // default is 60s
)

type CustomResponseWriter struct {
	writer *bufio.Writer
	header http.Header
}

func NewCustomResponseWriter(conn net.Conn) *CustomResponseWriter {
	return &CustomResponseWriter{
		writer: bufio.NewWriter(conn),
		header: http.Header{},
	}
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.header
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	nn, err := w.writer.Write(b)
	if err != nil {
		return 0, err
	}
	w.writer.Flush()
	return nn, nil
}
func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.writer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode)))
	w.header.Write(w.writer)
	w.writer.WriteString("\r\n")
	w.writer.Flush()
}

func handleRequest(r *http.Request, conn net.Conn, waitSig chan int) {
	w := NewCustomResponseWriter(conn)
	respondAsPlayground(w, r)

	waitSig <- 1
}

//// start child process
func StartWorker() error {
	pid := os.Getpid()
	lf := os.NewFile(uintptr(3), fmt.Sprintf("listener-%d", pid))
	defer lf.Close()

	pipeWriter := os.NewFile(uintptr(4), fmt.Sprintf("pipe-%d", pid))
	defer pipeWriter.Close()

	lc, err := net.FileListener(lf)
	if err != nil {
		return err
	}
	// get execution timeout
	timeout := defaultTimeout
	if t, err := strconv.Atoi(os.Getenv(EnvExecTimeout)); err == nil {
		timeout = t
	}

	// kill child process when parent process exits
	go watchMaster()

	for {
		conn, err := lc.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %v", err)
		}

		// set busy state
		writeProcState(pipeWriter, WORKER_STATE_BUSY)

		// Wrap the connection in a bufio.Reader to read the HTTP request
		bufReader := bufio.NewReader(conn)
		req, err := http.ReadRequest(bufReader)
		if err != nil {
			return fmt.Errorf("read HTTP request error: %v", err)
		}

		waitSig := make(chan int)
		go handleRequest(req, conn, waitSig)

		select {
		case <-waitSig:
			writeProcState(pipeWriter, WORKER_STATE_IDLE)
		case <-time.After(time.Duration(timeout) * time.Second):
			writeProcState(pipeWriter, WORKER_STATE_STOPPED)
			// wait for a while before exiting the process
			time.Sleep(100 * time.Millisecond)
			// exit the process directly
			os.Exit(1)
		}

		conn.Close()
	}
}

func watchMaster() {
	// if it is equal to 1 (init process ID),
	// it indicates that the master process has exited
	const watchInterval = 500 * time.Millisecond
	for range time.NewTicker(watchInterval).C {
		if os.Getppid() == 1 {
			os.Exit(1) //nolint:revive // Calling os.Exit is fine here in the prefork
		}
	}
}

func writeProcState(writer *os.File, state uint8) {
	pid := os.Getpid()
	buf := make([]byte, 5)
	// write state msg:
	// +-------------------------------------------+
	// | pid[3] | pid[2] | pid[1] | pid[0] | state |
	// +-------------------------------------------+

	binary.BigEndian.PutUint32(buf, uint32(pid))
	buf[4] = state

	// write state msg to notify parent process
	writer.Write(buf)
}

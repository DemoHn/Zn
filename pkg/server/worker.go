package server

import (
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

type FCGIHandler struct {
	pipeWriter *os.File
	timeout    int
}

func (f *FCGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	waitSig := make(chan int)

	go handleRequest(w, r, waitSig)

	select {
	case <-waitSig:
		writeProcState(f.pipeWriter, WORKER_STATE_IDLE)
	case <-time.After(time.Duration(f.timeout) * time.Second):
		writeProcState(f.pipeWriter, WORKER_STATE_STOPPED)
		// wait for a while before exiting the process
		time.Sleep(100 * time.Millisecond)
		// exit the process directly
		os.Exit(1)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request, waitSig chan int) {
	cgiEnvs := ProcessFCGIEnv(r)
	// read ZINC_ADAPTER
	if adapter, ok := cgiEnvs[EnvZincAdapter]; ok {
		switch adapter {
		case ZincAP_Playground:
			respondAsPlayground(w, r)
		case ZincAP_HTTPHandler:
			w.WriteHeader(200)
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte("TBD - http_handler"))
		}

		waitSig <- 1
		return
	}
	//// otherwise - return 403 directly

	w.WriteHeader(403)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("Invalid ZINC_ADAPTER"))

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
		// do actual request
		AcceptFCGIRequest(conn, &FCGIHandler{
			timeout:    timeout,
			pipeWriter: pipeWriter,
		})
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

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
	WORKER_STATE_INIT    uint8 = 1 << 0
	WORKER_STATE_IDLE    uint8 = 1 << 1
	WORKER_STATE_BUSY    uint8 = 1 << 2
	WORKER_STATE_STOPPED uint8 = 1 << 3
	WORKER_STATE_TIMEOUT uint8 = 1 << 4

	defaultTimeout = 60 // default is 60s
)

type FCGIHandler struct{}

func (f *FCGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cgiEnvs := ProcessFCGIEnv(r)
	// read ZINC_ADAPTER
	if adapter, ok := cgiEnvs[EnvZincAdapter]; ok {
		switch adapter {
		case ZincAP_Playground:
			respondAsPlayground(w, r)
			return
		case ZincAP_HTTPHandler:
			w.WriteHeader(200)
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte("TBD - http_handler"))
			return
		}
	}
	//// otherwise - return 403 directly

	w.WriteHeader(403)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("Invalid ZINC_ADAPTER"))
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
	writeProcState(pipeWriter, WORKER_STATE_IDLE)

	for {
		conn, err := lc.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %v", err)
		}

		acceptRequest(conn, pipeWriter, timeout)
	}
}

func acceptRequest(conn net.Conn, pipeWriter *os.File, timeout int) {
	waitSig := make(chan int)

	// set busy state
	writeProcState(pipeWriter, WORKER_STATE_BUSY)

	// do actual request
	go func() {
		AcceptFCGIRequest(conn, &FCGIHandler{})
		waitSig <- 1
	}()

	select {
	case <-waitSig:
		writeProcState(pipeWriter, WORKER_STATE_IDLE)
	case <-time.After(time.Duration(timeout) * time.Second):
		writeProcState(pipeWriter, WORKER_STATE_TIMEOUT)
		// wait for a while before exiting the process
		time.Sleep(100 * time.Millisecond)
		// when the work timeout, kill the worker process directly
		os.Exit(1)
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

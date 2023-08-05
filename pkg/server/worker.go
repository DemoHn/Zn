package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	WORKER_STATE_INIT    uint8 = 1 << 0
	WORKER_STATE_IDLE    uint8 = 1 << 1
	WORKER_STATE_BUSY    uint8 = 1 << 2
	WORKER_STATE_STOPPED uint8 = 1 << 3
	WORKER_STATE_TIMEOUT uint8 = 1 << 4
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

	// kill child process when parent process exits
	go watchMaster()

	writeProcState(pipeWriter, WORKER_STATE_IDLE)
	for {
		conn, err := lc.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %v", err)
		}
		AcceptFCGIRequest(conn, &FCGIHandler{})
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
	buf := make([]byte, 4)
	// write state msg:
	// +----------------------------------+
	// | pid[2] | pid[1] | pid[0] | state |
	// +----------------------------------+

	buf[0] = uint8(0xff & (pid >> 16))
	buf[1] = uint8(0xff & (pid >> 8))
	buf[2] = uint8(0xff & pid)
	buf[3] = state

	// write state msg to notify parent process
	writer.Write(buf)
}

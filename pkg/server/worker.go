package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
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

	fd4 := os.NewFile(uintptr(4), fmt.Sprintf("pipe-%d", pid))
	defer fd4.Close()

	lc, err := net.FileListener(lf)
	if err != nil {
		return err
	}

	// kill child process when parent process exits
	go watchMaster()

	go scheduleWriteTest(fd4)

	for {
		conn, err := lc.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %v", err)
		}
		log.Printf("当前处理请求PID = %d", pid)

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

func scheduleWriteTest(writer *os.File) {
	i := 0

	pid := os.Getpid()
	for {
		binary.Write(writer, binary.BigEndian, int32(pid))
		i++
		time.Sleep(time.Second)
	}
}

package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

type FCGIHandler struct{}

func (f *FCGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cgiEnvs := ProcessFCGIEnv(r)
	// read ZINC_ADAPTER
	if adapter, ok := cgiEnvs[EnvZincAdapter]; ok {
		switch adapter {
		case ZincAP_Playground:
			RespondAsPlayground(w, r)
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

	lc, err := net.FileListener(lf)
	if err != nil {
		return err
	}

	for {
		conn, err := lc.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %v", err)
		}

		AcceptFCGIRequest(conn, &FCGIHandler{})
	}
}

package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Zn-Server",
		Short: "Zn FastCGI 服务器",
		Long:  "Zn FastCGI 服务器",
		Run: func(c *cobra.Command, args []string) {
			StartServer()
		},
	}
)

type home struct{}

func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println(fcgi.ProcessEnv(r))
	fmt.Println(r.Header)
	fmt.Println(r.URL.Path, r.URL.Query())

	buf, _ := io.ReadAll(r.Body)

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<h2>This is HTTP Text</h2> body: <h3>%s</h3>", string(buf))))
}

func StartServer() {
	defer func() {
		os.Remove("/Users/demohn/test.sock")
	}()

	l, e := net.Listen("unix", "/Users/demohn/test.sock")
	if e != nil {
		fmt.Printf("error:%s", e)
	}

	fmt.Println("going to serve...")
	fcgi.Serve(l, &home{})
	fmt.Println("serve done...")
}

func main() {
	rootCmd.Execute()
}

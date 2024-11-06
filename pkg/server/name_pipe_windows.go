package server

import (
	"os"
	"syscall"
)

const gPipeNameFmt = "/tmp/zinc-server-pipe-%s"

type pipe struct {
	id string
}

func NewPipe(id string) *pipe {
	return &pipe{id: id}
}

func GetPipeID(p *pipe) string {
	return p.id
}

func CreateNamedPipe() (*pipe, error) {
	panic("Not Supported in Windows!")
}

func OpenNamedPipeReader(p *pipe) (*os.File, error) {
	panic("Not Supported in Windows!")
}

func ReadDataFromNamedPipe(pipeReader *os.File, b []byte) error {
	panic("Not Supported in Windows!")
}

func OpenNamedPipeWriter(p *pipe) (*os.File, error) {
	panic("Not Supported in Windows!")
}

func WriteDataToNamedPipe(pipeWriter *os.File, b []byte) error {
	panic("Not Supported in Windows!")
}

func CloseFD(f *os.File) {
	syscall.CloseOnExec(syscall.Handle(f.Fd()))
}

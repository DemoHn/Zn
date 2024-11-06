package server

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"syscall"
	"time"
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
	id := (rand.Int63() >> 32) | (time.Now().Unix() << 32)
	idStr := strconv.FormatInt(id, 16)

	if err := syscall.Mkfifo(fmt.Sprintf(gPipeNameFmt, idStr), 0666); err != nil {
		return nil, err
	}

	return &pipe{id: idStr}, nil
}

func OpenNamedPipeReader(p *pipe) (*os.File, error) {
	pipeFile := fmt.Sprintf(gPipeNameFmt, p.id)
	pipeReader, err := os.OpenFile(pipeFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}
	// os.Remove(pipeFile)
	return pipeReader, nil
}

func ReadDataFromNamedPipe(pipeReader *os.File, b []byte) error {
	_, err := pipeReader.Read(b)
	return err
}

func OpenNamedPipeWriter(p *pipe) (*os.File, error) {
	pipeFile := fmt.Sprintf(gPipeNameFmt, p.id)
	pipeWriter, err := os.OpenFile(pipeFile, os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	return pipeWriter, nil
}

func WriteDataToNamedPipe(pipeWriter *os.File, b []byte) error {
	_, err := pipeWriter.Write(b)
	return err
}

func CloseFD(f *os.File) {
	syscall.CloseOnExec(int(f.Fd()))
}

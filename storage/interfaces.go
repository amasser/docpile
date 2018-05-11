package storage

import (
	"io"
	"os"
)

type Reader interface {
	Read(string) (io.ReadCloser, error)
}
type Writer interface {
	Write(string, io.ReadCloser) error
}
type ReadWriter interface {
	Reader
	Writer
}

var (
	NotFoundError = os.ErrNotExist
)

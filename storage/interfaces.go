package storage

import (
	"io"
	"os"
)

type Reader interface {
	Read(string) (io.ReadCloser, error)
}

// NOTE: writer works for create and append, it depends upon the underlying implementation.
type Writer interface {
	Write(string, io.ReadCloser) error
}

var NotFoundError = os.ErrNotExist

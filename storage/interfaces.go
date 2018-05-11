package storage

import (
	"io"
	"os"
	"reflect"
)

type Reader interface {
	Read(string) (io.ReadCloser, error)
}

// NOTE: writer works for create and append, it depends upon the underlying implementation.
type Writer interface {
	Write(string, io.ReadCloser) error
}

type ReadWriter interface {
	Reader
	Writer
}

var NotFoundError = os.ErrNotExist

type Registry interface {
	Name(reflect.Type) (string, error)
	Type(string) (reflect.Type, error)
}

type Serializer interface {
	Serialize(interface{}, io.Writer) error
	Deserialize(io.Reader, interface{}) error
}

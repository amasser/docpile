package storage

import (
	"io"
	"os"
	"path"
	"strings"
)

type FileStorage struct {
	workspace  string
	writeFlags int
}

func NewFileStorage(workspace string, options ...FileOption) *FileStorage {
	if len(workspace) == 0 {
		panic("workspace is required")
	}

	this := &FileStorage{workspace: workspace, writeFlags: os.O_CREATE | os.O_WRONLY}

	for _, option := range options {
		option(this)
	}

	return this
}
func (this *FileStorage) composeFilename(key string) string {
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		panic("key is required")
	}
	return path.Join(this.workspace, key)
}

func (this *FileStorage) Read(key string) (io.ReadCloser, error) {
	key = this.composeFilename(key)
	if handle, err := os.Open(key); err == nil {
		return handle, nil
	} else if os.IsNotExist(err) {
		return nil, NotFoundError
	} else {
		return nil, err
	}
}

func (this *FileStorage) Write(key string, source io.ReadCloser) error {
	key = this.composeFilename(key)
	if handle, err := os.OpenFile(key, this.writeFlags, 0644); err == nil {
		return this.write(source, handle)
	} else {
		return err
	}
}
func (this *FileStorage) write(source io.ReadCloser, destination io.WriteCloser) error {
	defer source.Close()
	defer destination.Close()
	var buffer [1024 * 16]byte
	_, err := io.CopyBuffer(destination, source, buffer[:])
	return err
}

func (this *FileStorage) ensureWorkspace() {
	if err := os.MkdirAll(this.workspace, 0755); err != nil {
		panic(err)
	}
}

type FileOption func(*FileStorage)

func Append() FileOption          { return func(this *FileStorage) { this.writeFlags |= os.O_APPEND } }
func EnsureWorkspace() FileOption { return func(this *FileStorage) { this.ensureWorkspace() } }

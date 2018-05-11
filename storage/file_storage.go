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

func NewFileStorage(workspace string) *FileStorage {
	this := &FileStorage{
		workspace:  workspace,
		writeFlags: os.O_CREATE | os.O_WRONLY,
	}

	this.ensureWorkspace() // TODO: functional options

	return this
}
func (this *FileStorage) ensureWorkspace() {
	if len(this.workspace) == 0 {
		panic("workspace is required")
	}

	if err := os.MkdirAll(this.workspace, 0755); err != nil {
		panic(err)
	}
}
func (this *FileStorage) composeFilename(key string) string {
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		panic("key is required")
	}
	return path.Join(this.workspace, key)
}
func (this *FileStorage) Append() *FileStorage {
	this.writeFlags |= os.O_APPEND // TODO: functional options
	return this
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

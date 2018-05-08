package storage

import (
	"io"
	"os"
	"path"
	"strings"
)

type LocalStorage struct {
	workspace string
}

func NewLocalStorage(workspace string) *LocalStorage {
	if len(workspace) == 0 {
		panic("workspace is required")
	}

	return &LocalStorage{workspace: workspace}
}

func (this *LocalStorage) EnsureWorkspace() error {
	return os.MkdirAll(this.workspace, 0755)
}

func (this *LocalStorage) Read(key string) (io.ReadCloser, error) {
	key = this.composeFilename(key)
	if handle, err := os.Open(key); err == nil {
		return handle, nil
	} else if os.IsNotExist(err) {
		return nil, NotFoundError
	} else {
		return nil, err
	}
}

func (this *LocalStorage) Write(key string, source io.ReadCloser) error {
	key = this.composeFilename(key)
	if handle, err := os.OpenFile(key, os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		return this.write(source, handle)
	} else {
		return err
	}
}
func (this *LocalStorage) write(source io.ReadCloser, destination io.WriteCloser) error {
	defer source.Close()
	defer destination.Close()
	var buffer [1024 * 16]byte
	_, err := io.CopyBuffer(destination, source, buffer[:])
	return err
}

func (this *LocalStorage) Move(old, new string) error {
	old = this.composeFilename(old)
	new = this.composeFilename(new)
	if err := os.Rename(old, new); os.IsNotExist(err) {
		return NotFoundError
	} else if err != nil {
		return err
	} else {
		return nil
	}
}

func (this *LocalStorage) composeFilename(key string) string {
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		panic("key is required")
	}
	return path.Join(this.workspace, key)
}

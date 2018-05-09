package storage

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/events"
)

type LocalStorageHandler struct {
	inner  domain.Handler
	writer Writer
}

func NewLocalStorageHandler(inner domain.Handler, writer Writer) *LocalStorageHandler {
	return &LocalStorageHandler{inner: inner, writer: writer}
}

func (this *LocalStorageHandler) Handle(message interface{}) (uint64, error) {
	switch message := message.(type) {
	case domain.ImportManagedStreamingAsset:
		return this.handleImportManagedStreamingAsset(message)
	default:
		return this.inner.Handle(message)
	}
}
func (this *LocalStorageHandler) handleImportManagedStreamingAsset(message domain.ImportManagedStreamingAsset) (uint64, error) {
	buffer, err := bufferStream(message.Size, message.Body)
	if err != nil {
		return 0, err
	}

	id, err := this.sendMessage(message.Name, message.MIMEType, buffer)
	if err != nil {
		return 0, err
	}

	if err := this.writeBuffer(message.Name, id, buffer); err != nil {
		panic(err) // because we don't have a compensating event, we panic so the event never happens
	}

	return id, nil
}

func (this *LocalStorageHandler) sendMessage(name, mime string, buffer *bytes.Buffer) (uint64, error) {
	return this.Handle(domain.ImportManagedAsset{
		Name:     name,
		MIMEType: mime,
		Hash:     computeHash(buffer),
	})
}
func (this *LocalStorageHandler) writeBuffer(originalFilename string, id uint64, buffer *bytes.Buffer) error {
	filename := fmt.Sprintf("%d%s", id, path.Ext(originalFilename))
	source := ioutil.NopCloser(buffer)
	return this.writer.Write(filename, source)
}

func bufferStream(size uint64, reader io.ReadCloser) (*bytes.Buffer, error) {
	defer reader.Close()
	buffer := bytes.NewBuffer(make([]byte, int(size)))
	if _, err := io.Copy(buffer, reader); err == nil {
		return buffer, nil
	} else {
		return nil, err
	}
}
func computeHash(buffer *bytes.Buffer) events.SHA256Hash {
	return events.SHA256Hash(sha256.Sum256(buffer.Bytes()))
}

package domain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"github.com/joliver/docpile/app/events"
	"github.com/joliver/docpile/generic/handlers"
	"github.com/joliver/docpile/generic/storage"
)

type WriteAssetHandler struct {
	inner  handlers.Handler
	writer storage.Writer
}

func NewWriteAssetHandler(inner handlers.Handler, writer storage.Writer) *WriteAssetHandler {
	return &WriteAssetHandler{inner: inner, writer: writer}
}

func (this *WriteAssetHandler) Handle(message interface{}) handlers.Result {
	switch message := message.(type) {
	case ImportManagedStreamingAsset:
		return this.handle(message)
	default:
		return this.inner.Handle(message)
	}
}
func (this *WriteAssetHandler) handle(message ImportManagedStreamingAsset) handlers.Result {
	buffer, err := bufferStream(message.Size, message.Body)
	if err != nil {
		return newResult(0, err)
	}

	result := this.sendMessage(message.Name, message.MIMEType, buffer)
	if result.Error != nil {
		return result
	}

	if err := this.writeBuffer(result.ID, message.Name, buffer); err != nil {
		panic(err) // because we don't have a compensating event, we panic so the event never happens
	}

	return result
}

func (this *WriteAssetHandler) sendMessage(name, mime string, buffer *bytes.Buffer) handlers.Result {
	return this.inner.Handle(ImportManagedAsset{
		Name:     name,
		MIMEType: mime,
		Hash:     computeHash(buffer),
	})
}
func (this *WriteAssetHandler) writeBuffer(id uint64, name string, buffer *bytes.Buffer) error {
	filename := fmt.Sprintf("%d%s", id, path.Ext(name))
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

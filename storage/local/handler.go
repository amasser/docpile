package local

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/events"
	"bitbucket.org/jonathanoliver/docpile/storage"
)

type Handler struct {
	inner  domain.Handler
	writer storage.Writer
}

func NewHandler(inner domain.Handler, writer storage.Writer) *Handler {
	return &Handler{inner: inner, writer: writer}
}

func (this *Handler) Handle(message interface{}) (uint64, error) {
	switch message := message.(type) {
	case domain.ImportManagedStreamingAsset:
		return this.handle(message)
	default:
		return this.inner.Handle(message)
	}
}
func (this *Handler) handle(message domain.ImportManagedStreamingAsset) (uint64, error) {
	buffer, err := bufferStream(message.Size, message.Body)
	if err != nil {
		return 0, err
	}

	id, err := this.sendMessage(message.Name, message.MIMEType, buffer)
	if err != nil {
		return 0, err
	}

	if err := this.writeBuffer(id, message.Name, buffer); err != nil {
		panic(err) // because we don't have a compensating event, we panic so the event never happens
	}

	return id, nil
}

func (this *Handler) sendMessage(name, mime string, buffer *bytes.Buffer) (uint64, error) {
	return this.inner.Handle(domain.ImportManagedAsset{
		Name:     name,
		MIMEType: mime,
		Hash:     computeHash(buffer),
	})
}
func (this *Handler) writeBuffer(id uint64, name string, buffer *bytes.Buffer) error {
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

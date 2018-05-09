package storage

import "bitbucket.org/jonathanoliver/docpile/domain"

type LocalStorageHandler struct {
	inner     domain.Handler
	writer    Writer
	workspace string
}

func NewLocalStorageHandler(inner domain.Handler, writer Writer, workspace string) *LocalStorageHandler {
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
	return 0, nil
	//// get temp filename
	//// write to that location and compute hash while doing so
	//// pipe through domain using new command
	//// move file to permanent location
	//// if fail, send DeleteManagedAsset
	//return this.inner.Handle(domain.ImportManagedAsset{
	//	Name:     message.Name,
	//	MIMEType: message.MIMEType,
	//})
}

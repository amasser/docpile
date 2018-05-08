package storage

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
)

type LocalStorageHandler struct {
	inner   domain.Handler
	storage Writer
}

func NewLocalStorageHandler(inner domain.Handler, storage Writer) *LocalStorageHandler {
	return &LocalStorageHandler{inner: inner, storage: storage}
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
	// get temp filename
	// write to that location and compute hash while doing so
	// pipe through domain using new command
	// move file to permanent location
	// if fail, send DeleteManagedAsset
	return 0, nil
}

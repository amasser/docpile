package domain

import (
	"fmt"
	"reflect"
	"log"
)

type Handler struct {
	root       *Aggregate
	applicator MessageApplicator
}

func NewHandler(app *Aggregate) *Handler {
	return &Handler{root: app}
}

func (this *Handler) Handle(message interface{}) (uint64, error) {
	if id, err := this.handle(message); err != nil {
		return 0, err
	} else {
		this.applicator.Apply(this.root.Consume())
		return id, nil
	}
}
func (this *Handler) handle(message interface{}) (uint64, error) {
	switch message := message.(type) {
	case AddTag:
		return this.root.AddTag(message.Name)
	case ImportManagedAsset:
		return this.root.ImportManagedAsset(message.Name, message.MIMEType, message.Hash)
	case ImportCloudAsset:
		return this.root.ImportCloudAsset(message.Name, message.Provider, message.Resource)
	case DefineDocument:
		return this.root.DefineDocument(message.Document)
	default:
		log.Panicf(fmt.Sprintf("Handler cannot handle '%s'", reflect.TypeOf(message)))
	}

	return 0, nil
}

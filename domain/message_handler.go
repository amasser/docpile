package domain

import (
	"fmt"
	"log"
	"reflect"
)

type MessageHandler struct {
	root       *Aggregate
	store      EventStore
	applicator Applicator
}

func NewMessageHandler(root *Aggregate, store EventStore, applicator Applicator) *MessageHandler {
	return &MessageHandler{root: root, store: store, applicator: applicator}
}

func (this *MessageHandler) Handle(message interface{}) (uint64, error) {
	if id, err := this.handle(message); err != nil {
		return 0, err
	} else if messages := this.root.Consume(); len(messages) == 0 {
		return id, nil
	} else if err = this.store.Store(messages); err != nil {
		panic(err)
	} else {
		this.applicator.Apply(messages)
		return id, nil
	}
}
func (this *MessageHandler) handle(message interface{}) (uint64, error) {
	switch message := message.(type) {

	case AddTag:
		return this.root.AddTag(message.Name)
	case RenameTag:
		return this.root.RenameTag(message.ID, message.Name)
	case DefineTagSynonym:
		return this.root.DefineTagSynonym(message.ID, message.Name)
	case RemoveTagSynonym:
		return this.root.RemoveTagSynonym(message.ID, message.Name)

	case ImportManagedAsset:
		return this.root.ImportManagedAsset(message.Name, message.MIMEType, message.Hash)
	case ImportCloudAsset:
		return this.root.ImportCloudAsset(message.Name, message.Provider, message.Resource)

	case DefineDocument:
		return this.root.DefineDocument(message.Document)

	default:
		log.Panicf(fmt.Sprintf("MessageHandler cannot handle '%s'", reflect.TypeOf(message)))
	}

	return 0, nil
}

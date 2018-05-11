package domain

import (
	"fmt"
	"log"
	"reflect"

	"bitbucket.org/jonathanoliver/docpile/infrastructure"
)

type CommandHandler struct {
	root       *Aggregate
	applicator infrastructure.Applicator
}

func NewCommandHandler(root *Aggregate, applicator infrastructure.Applicator) *CommandHandler {
	return &CommandHandler{root: root, applicator: applicator}
}

func (this *CommandHandler) Handle(message interface{}) infrastructure.Result {
	result := this.handle(message)
	if result.Error == nil {
		this.applicator.Apply(this.root.Consume())
	}

	return result
}
func (this *CommandHandler) handle(message interface{}) infrastructure.Result {
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
		log.Panicf(fmt.Sprintf("CommandHandler cannot handle '%s'", reflect.TypeOf(message)))
		return newResult(0, nil)
	}
}

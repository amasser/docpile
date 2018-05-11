package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/infrastructure"
	"github.com/smartystreets/detour"
)

type DocumentWriteController struct {
	handler infrastructure.Handler
}

func NewDocumentWriteController(handler infrastructure.Handler) *DocumentWriteController {
	return &DocumentWriteController{handler: handler}
}

func (this *DocumentWriteController) Define(input *inputs.DefineDocument) detour.Renderer {
	if documentID, err := this.define(input); err == nil {
		return newEntityResult(documentID)
	} else if err == domain.AssetNotFoundError {
		return inputs.AssetDoesNotExistResult
	} else if err == domain.TagNotFoundError {
		return inputs.TagDoesNotExistResult
	} else if err == domain.DocumentNotFoundError {
		return inputs.DocumentDoesNotExistResult
	} else {
		return UnknownErrorResult
	}
}
func (this *DocumentWriteController) define(input *inputs.DefineDocument) (uint64, error) {
	return this.handler.Handle(domain.DefineDocument{
		Document: domain.DocumentDefinition{
			AssetID:     input.AssetID,
			AssetOffset: input.AssetOffset,
			Published:   input.Published,
			PeriodBegin: input.PeriodBegin,
			PeriodEnd:   input.PeriodEnd,
			Tags:        input.Tags,
			Documents:   input.Documents,
			Description: input.Description,
		},
	})

}

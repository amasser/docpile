package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"github.com/smartystreets/detour"
)

type DocumentController struct {
	handler domain.Handler
}

func NewDocumentController(handler domain.Handler) *DocumentController {
	return &DocumentController{handler: handler}
}

func (this *DocumentController) Define(input *inputs.DefineDocument) detour.Renderer {
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
func (this *DocumentController) define(input *inputs.DefineDocument) (uint64, error) {
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

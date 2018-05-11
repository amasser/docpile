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
	if result := this.define(input); result.Error == nil {
		return newEntityResult(result.ID)
	} else if result.Error == domain.AssetNotFoundError {
		return inputs.AssetDoesNotExistResult
	} else if result.Error == domain.TagNotFoundError {
		return inputs.TagDoesNotExistResult
	} else if result.Error == domain.DocumentNotFoundError {
		return inputs.DocumentDoesNotExistResult
	} else {
		return UnknownErrorResult
	}
}
func (this *DocumentWriteController) define(input *inputs.DefineDocument) infrastructure.Result {
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

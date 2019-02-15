package http

import (
	"github.com/joliver/docpile/app/domain"
	"github.com/joliver/docpile/app/http/inputs"
	"github.com/joliver/docpile/generic/handlers"
	"github.com/smartystreets/detour"
)

type DocumentWriter struct {
	handler handlers.Handler
}

func NewDocumentWriter(handler handlers.Handler) *DocumentWriter {
	return &DocumentWriter{handler: handler}
}

func (this *DocumentWriter) Define(input *inputs.DefineDocument) detour.Renderer {
	return this.renderResult(domain.DefineDocument{
		Document: domain.DocumentDefinition{
			AssetID:     input.AssetID,
			AssetOffset: input.AssetOffset,
			Published:   input.Published,
			PeriodMin:   input.PeriodMin,
			PeriodMax:   input.PeriodMax,
			Tags:        input.Tags,
			Documents:   input.Documents,
			Description: input.Description,
		},
	})
}
func (this *DocumentWriter) renderResult(message interface{}) detour.Renderer {
	if result := this.handler.Handle(message); result.Error == nil {
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

func (this *DocumentWriter) Remove(input *inputs.IDInput) detour.Renderer {
	message := domain.RemoveDocument{ID: input.ID}
	result := this.handler.Handle(message)
	if result.Error == domain.DocumentNotFoundError {
		return inputs.IDNotFoundResult
	} else if result.Error != nil {
		return UnknownErrorResult
	} else {
		return nil // OK
	}
}

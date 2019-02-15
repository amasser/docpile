package http

import (
	"github.com/joliver/docpile/app/domain"
	"github.com/joliver/docpile/app/http/inputs"
	"github.com/joliver/docpile/generic/handlers"
	"github.com/smartystreets/detour"
)

type AssetWriter struct {
	handler handlers.Handler
}

func NewAssetWriter(handler handlers.Handler) *AssetWriter {
	return &AssetWriter{handler: handler}
}

func (this *AssetWriter) ImportManaged(input *inputs.ImportManagedAsset) detour.Renderer {
	if result := this.importManaged(input); result.Error == nil {
		return newEntityResult(result.ID)
	} else if result.Error == domain.AssetAlreadyExistsError {
		return inputs.DuplicateAssetResult
	} else if result.Error == domain.StoreAssetError {
		return UnknownErrorResult
	} else {
		return UnknownErrorResult
	}
}
func (this *AssetWriter) importManaged(input *inputs.ImportManagedAsset) handlers.Result {
	return this.handler.Handle(domain.ImportManagedStreamingAsset{
		Name:     input.Name,
		MIMEType: input.MIMEType,
		Size:     input.Size,
		Body:     input.Reader,
	})
}

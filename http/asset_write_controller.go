package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/library"
	"github.com/smartystreets/detour"
)

type AssetWriteController struct {
	handler library.Handler
}

func NewAssetWriteController(handler library.Handler) *AssetWriteController {
	return &AssetWriteController{handler: handler}
}

func (this *AssetWriteController) ImportManaged(input *inputs.ImportManagedAsset) detour.Renderer {
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
func (this *AssetWriteController) importManaged(input *inputs.ImportManagedAsset) library.Result {
	return this.handler.Handle(domain.ImportManagedStreamingAsset{
		Name:     input.Name,
		MIMEType: input.MIMEType,
		Size:     input.Size,
		Body:     input.Reader,
	})
}

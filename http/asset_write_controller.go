package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/infrastructure"
	"github.com/smartystreets/detour"
)

type AssetWriteController struct {
	handler infrastructure.Handler
}

func NewAssetWriteController(handler infrastructure.Handler) *AssetWriteController {
	return &AssetWriteController{handler: handler}
}

func (this *AssetWriteController) ImportManaged(input *inputs.ImportManagedAsset) detour.Renderer {
	if assetID, err := this.importManaged(input); err == nil {
		return newEntityResult(assetID)
	} else if err == domain.AssetAlreadyExistsError {
		return inputs.DuplicateAssetResult
	} else if err == domain.StoreAssetError {
		return UnknownErrorResult
	} else {
		return UnknownErrorResult
	}
}
func (this *AssetWriteController) importManaged(input *inputs.ImportManagedAsset) (uint64, error) {
	return this.handler.Handle(domain.ImportManagedStreamingAsset{
		Name:     input.Name,
		MIMEType: input.MIMEType,
		Size:     input.Size,
		Body:     input.Reader,
	})
}

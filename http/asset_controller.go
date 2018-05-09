package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"github.com/smartystreets/detour"
)

type AssetController struct {
	handler domain.Handler
}

func NewAssetController(handler domain.Handler) *AssetController {
	return &AssetController{handler: handler}
}

func (this *AssetController) Import(input *inputs.ImportManagedAsset) detour.Renderer {
	if assetID, err := this.importAsset(input); err == nil {
		return newEntityResult(assetID)
	} else if err == domain.AssetAlreadyExistsError {
		return inputs.DuplicateAssetResult
	} else if err == domain.StoreAssetError {
		return UnknownErrorResult
	} else {
		return UnknownErrorResult
	}
}
func (this *AssetController) importAsset(input *inputs.ImportManagedAsset) (uint64, error) {
	return this.handler.Handle(domain.ImportManagedStreamingAsset{
		Name:     input.Name,
		MIMEType: input.MIMEType,
		Size:     input.Size,
		Body:     input.Reader,
	})
}

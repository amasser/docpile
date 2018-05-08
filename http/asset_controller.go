package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"github.com/smartystreets/detour"
)

type AssetController struct {
	app domain.ManagedAssetStreamImporter
}

func NewAssetController(app domain.ManagedAssetStreamImporter) *AssetController {
	return &AssetController{app: app}
}

func (this *AssetController) Add(input *inputs.ImportManagedAsset) detour.Renderer {
	if assetID, err := this.app.ImportManagedAsset(input.Name, input.MIMEType, input.Reader); err == nil {
		return newEntityResult(assetID)
	} else if err == domain.AssetAlreadyExistsError {
		return inputs.DuplicateAssetResult
	} else if err == domain.StoreAssetError {
		return UnknownErrorResult
	} else {
		return UnknownErrorResult
	}
}

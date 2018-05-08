package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"github.com/smartystreets/detour"
)

type DocumentController struct {
	app domain.DocumentDefiner
}

func NewDocumentController(app domain.DocumentDefiner) *DocumentController {
	return &DocumentController{app: app}
}

func (this *DocumentController) Add(input *inputs.DefineDocument) detour.Renderer {
	proposed := domain.DocumentDefinition{
		AssetID:     input.AssetID,
		AssetOffset: input.AssetOffset,
		Published:   input.Published,
		PeriodBegin: input.PeriodBegin,
		PeriodEnd:   input.PeriodEnd,
		Tags:        input.Tags,
		Documents:   input.Documents,
		Description: input.Description,
	}

	if documentID, err := this.app.DefineDocument(proposed); err == nil {
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

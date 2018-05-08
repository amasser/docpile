package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"github.com/smartystreets/detour"
)

type TagController struct {
	app domain.TagAdder
}

func NewTagController(app domain.TagAdder) *TagController {
	return &TagController{app: app}
}

func (this *TagController) Add(input *inputs.AddTag) detour.Renderer {
	if tagID, err := this.app.AddTag(input.Name); err == nil {
		return newEntityResult(tagID)
	} else if err == domain.TagAlreadyExistsError {
		return inputs.DuplicateTagResult
	} else {
		return UnknownErrorResult
	}
}

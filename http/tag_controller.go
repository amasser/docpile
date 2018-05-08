package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"github.com/smartystreets/detour"
)

type TagController struct {
	handler domain.Handler
}

func NewTagController(handler domain.Handler) *TagController {
	return &TagController{handler: handler}
}

func (this *TagController) Add(input *inputs.AddTag) detour.Renderer {
	if tagID, err := this.add(input); err == nil {
		return newEntityResult(tagID)
	} else if err == domain.TagAlreadyExistsError {
		return inputs.DuplicateTagResult
	} else {
		return UnknownErrorResult
	}
}
func (this *TagController) add(input *inputs.AddTag) (uint64, error) {
	return this.handler.Handle(domain.AddTag{Name: input.Name})
}

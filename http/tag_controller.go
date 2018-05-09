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

func (this *TagController) Rename(input *inputs.RenameTag) detour.Renderer {
	if err := this.rename(input); err == nil {
		return nil
	} else if err == domain.TagAlreadyExistsError {
		return inputs.DuplicateTagResult
	} else if err == domain.TagNotFoundError {
		return inputs.TagNotFoundResult
	} else {
		return UnknownErrorResult
	}
}
func (this *TagController) rename(input *inputs.RenameTag) error {
	_, err := this.handler.Handle(domain.RenameTag{
		ID:   input.ID,
		Name: input.Name,
	})
	return err
}

func (this *TagController) DefineSynonym(input *inputs.DefineTagSynonym) detour.Renderer {
	if err := this.defineSynonym(input); err == nil {
		return nil
	} else if err == domain.TagAlreadyExistsError {
		return inputs.DuplicateTagResult
	} else if err == domain.TagNotFoundError {
		return inputs.TagNotFoundResult
	} else {
		return UnknownErrorResult
	}
}
func (this *TagController) defineSynonym(input *inputs.DefineTagSynonym) error {
	_, err := this.handler.Handle(domain.DefineTagSynonym{
		ID:   input.ID,
		Name: input.Name,
	})
	return err
}

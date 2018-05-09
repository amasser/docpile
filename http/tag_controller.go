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
	return this.renderTagResult(domain.AddTag{Name: input.Name})
}
func (this *TagController) Rename(input *inputs.RenameTag) detour.Renderer {
	return this.renderTagResult(domain.RenameTag{ID: input.ID, Name: input.Name})
}

func (this *TagController) DefineSynonym(input *inputs.DefineTagSynonym) detour.Renderer {
	return this.renderTagResult(domain.DefineTagSynonym{ID: input.ID, Name: input.Name})
}
func (this *TagController) RemoveSynonym(input *inputs.RemoveTagSynonym) detour.Renderer {
	return this.renderTagResult(domain.RemoveTagSynonym{ID: input.ID, Name: input.Name})
}

func (this *TagController) renderTagResult(message interface{}) detour.Renderer {
	if id, err := this.handler.Handle(message); id == 0 && err == nil {
		return nil
	} else if id > 0 && err == nil {
		return newEntityResult(id)
	} else if err == domain.TagAlreadyExistsError {
		return inputs.DuplicateTagResult
	} else if err == domain.TagNotFoundError {
		return inputs.TagNotFoundResult
	} else {
		return UnknownErrorResult
	}
}

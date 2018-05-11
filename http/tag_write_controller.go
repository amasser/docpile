package http

import (
	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/infrastructure"
	"github.com/smartystreets/detour"
)

type TagWriteController struct {
	handler infrastructure.Handler
}

func NewTagWriteController(handler infrastructure.Handler) *TagWriteController {
	return &TagWriteController{handler: handler}
}

func (this *TagWriteController) Add(input *inputs.AddTag) detour.Renderer {
	return this.renderTagResult(domain.AddTag{Name: input.Name})
}
func (this *TagWriteController) Rename(input *inputs.RenameTag) detour.Renderer {
	return this.renderTagResult(domain.RenameTag{ID: input.ID, Name: input.Name})
}

func (this *TagWriteController) DefineSynonym(input *inputs.DefineTagSynonym) detour.Renderer {
	return this.renderTagResult(domain.DefineTagSynonym{ID: input.ID, Name: input.Name})
}
func (this *TagWriteController) RemoveSynonym(input *inputs.RemoveTagSynonym) detour.Renderer {
	return this.renderTagResult(domain.RemoveTagSynonym{ID: input.ID, Name: input.Name})
}

func (this *TagWriteController) renderTagResult(message interface{}) detour.Renderer {
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

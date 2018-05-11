package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/domain"
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/library"
	"github.com/smartystreets/detour"
)

type TagWriteController struct {
	handler library.Handler
}

func NewTagWriteController(handler library.Handler) *TagWriteController {
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
	if result := this.handler.Handle(message); result.ID == 0 && result.Error == nil {
		return nil
	} else if result.ID > 0 && result.Error == nil {
		return newEntityResult(result.ID)
	} else if result.Error == domain.TagAlreadyExistsError {
		return inputs.DuplicateTagResult
	} else if result.Error == domain.TagNotFoundError {
		return inputs.TagNotFoundResult
	} else {
		return UnknownErrorResult
	}
}

package http

import (
	"github.com/joliver/docpile/app/domain"
	"github.com/joliver/docpile/app/http/inputs"
	"github.com/joliver/docpile/generic/handlers"
	"github.com/smartystreets/detour"
)

type TagWriter struct {
	handler handlers.Handler
}

func NewTagWriter(handler handlers.Handler) *TagWriter {
	return &TagWriter{handler: handler}
}

func (this *TagWriter) Add(input *inputs.AddTag) detour.Renderer {
	return this.renderTagResult(domain.AddTag{Name: input.Name})
}
func (this *TagWriter) Remove(input *inputs.IDInput) detour.Renderer {
	result := this.handler.Handle(domain.RemoveTag{ID: input.ID})
	if result.Error == domain.TagNotFoundError {
		return inputs.IDNotFoundResult
	} else if result.Error != nil {
		return UnknownErrorResult
	} else {
		return nil
	}
}
func (this *TagWriter) Rename(input *inputs.TagInput) detour.Renderer {
	return this.renderTagResult(domain.RenameTag{ID: input.ID, Name: input.Name})
}

func (this *TagWriter) DefineSynonym(input *inputs.TagInput) detour.Renderer {
	return this.renderTagResult(domain.DefineTagSynonym{ID: input.ID, Name: input.Name})
}
func (this *TagWriter) RemoveSynonym(input *inputs.TagInput) detour.Renderer {
	return this.renderTagResult(domain.RemoveTagSynonym{ID: input.ID, Name: input.Name})
}

func (this *TagWriter) renderTagResult(message interface{}) detour.Renderer {
	if result := this.handler.Handle(message); result.ID == 0 && result.Error == nil {
		return nil
	} else if result.ID > 0 && result.Error == nil {
		return newEntityResult(result.ID)
	} else if result.Error == domain.TagAlreadyExistsError {
		return inputs.DuplicateTagResult
	} else if result.Error == domain.TagNotFoundError {
		return inputs.TagNotFoundResult
	} else if result.Error == domain.SynonymNotFoundError {
		return inputs.SynonymNotFoundResult
	} else {
		return UnknownErrorResult
	}
}

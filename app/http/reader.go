package http

import (
	"github.com/joliver/docpile/app/http/inputs"
	"github.com/joliver/docpile/app/projections"
	"github.com/smartystreets/detour"
)

type Reader struct {
	tags      *projections.AllTags
	documents *projections.AllDocuments
}

func NewReader(tags *projections.AllTags, documents *projections.AllDocuments) *Reader {
	return &Reader{tags: tags, documents: documents}
}

func (this *Reader) ListTags() detour.Renderer {
	return jsonResult(this.tags.List())
}
func (this *Reader) ListDocuments() detour.Renderer {
	return jsonResult(this.documents.List())
}

func (this *Reader) LoadTag(input *inputs.IDInput) detour.Renderer {
	return this.render(this.tags.Load(input.ID))
}
func (this *Reader) LoadDocument(input *inputs.IDInput) detour.Renderer {
	return this.render(this.documents.Load(input.ID))
}
func (this *Reader) render(value interface{}, err error) detour.Renderer {
	if err != nil {
		return inputs.IDNotFoundResult
	} else {
		return jsonResult(value)
	}
}

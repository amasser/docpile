package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"github.com/smartystreets/detour"
)

type Reader struct {
	reader projectionReader
}

func NewReader(reader projectionReader) *Reader {
	return &Reader{reader: reader}
}

func (this *Reader) ListTags() detour.Renderer {
	return jsonResult(this.reader.ListTags())
}
func (this *Reader) ListDocuments() detour.Renderer {
	return jsonResult(this.reader.ListDocuments())
}

func (this *Reader) LoadTag(input *inputs.IDInput) detour.Renderer {
	return this.render(this.reader.LoadTag(input.ID))
}
func (this *Reader) LoadDocument(input *inputs.IDInput) detour.Renderer {
	return this.render(this.reader.LoadDocument(input.ID))
}
func (this *Reader) render(value interface{}, err error) detour.Renderer {
	if err != nil {
		return inputs.IDNotFoundResult
	} else {
		return jsonResult(value)
	}
}

type projectionReader interface {
	ListTags() interface{}
	LoadTag(uint64) (interface{}, error)
	ListDocuments() interface{}
	LoadDocument(uint64) (interface{}, error)
}

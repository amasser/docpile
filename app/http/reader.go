package http

import "github.com/smartystreets/detour"

type Reader struct {
	reader projectionReader
}

func NewReader(reader projectionReader) *Reader {
	return &Reader{reader: reader}
}

func (this *Reader) ListTags() detour.Renderer {
	return detour.JSONResult{Content: this.reader.ListTags()}
}
func (this *Reader) ListDocuments() detour.Renderer {
	return detour.JSONResult{Content: this.reader.ListDocuments()}
}

// TODO: load ID input model (read path)
type projectionReader interface {
	ListTags() interface{}
	LoadTag(uint64) (interface{}, error)
	ListDocuments() interface{}
	LoadDocument(uint64) (interface{}, error)
}

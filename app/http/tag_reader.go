package http

import "github.com/smartystreets/detour"

type TagReader struct {
	reader tagLister
}

func NewTagReader(reader tagLister) *TagReader {
	return &TagReader{reader: reader}
}

func (this *TagReader) List() detour.Renderer {
	return detour.JSONResult{Content: this.reader.ListTags()}
}

type tagLister interface {
	ListTags() interface{}
}

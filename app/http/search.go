package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/app/projections"
	"github.com/smartystreets/detour"
)

type Search struct {
	search searcher
}

func NewSearch(search searcher) *Search {
	return &Search{search: search}
}

func (this *Search) Documents(input *inputs.SearchDocument) detour.Renderer {
	spec := projections.NewDocumentSearch(
		input.PublishedMin, input.PublishedMax,
		input.PeriodMin, input.PeriodMax,
		input.Tags)
	return jsonResult(this.search.SearchDocuments(spec))
}
func (this *Search) Tags(input *inputs.SearchTag) detour.Renderer {
	return nil
}

type searcher interface {
	SearchDocuments(projections.DocumentSpecification) interface{}
}

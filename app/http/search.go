package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/app/search"
	"github.com/smartystreets/detour"
)

type Search struct {
	documents *search.DocumentSearcher
}

func NewSearch(documents *search.DocumentSearcher) *Search {
	return &Search{documents: documents}
}

func (this *Search) Documents(input *inputs.SearchDocument) detour.Renderer {
	return jsonResult(this.documents.Search(search.NewDocumentSpecification(
		input.PublishedMin, input.PublishedMax,
		input.PeriodMin, input.PeriodMax,
		input.Tags)))
}
func (this *Search) Tags(input *inputs.SearchTag) detour.Renderer {
	return nil
}

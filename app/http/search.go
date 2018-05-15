package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/app/projections"
	"github.com/smartystreets/detour"
)

type Search struct {
	documents *projections.AllDocuments
	tags      *projections.TagSearch
}

func NewSearch(documents *projections.AllDocuments, tags *projections.TagSearch) *Search {
	return &Search{documents: documents, tags: tags}
}

func (this *Search) Documents(input *inputs.SearchDocument) detour.Renderer {
	criteria := projections.NewDocumentCriteria(
		input.PublishedMin, input.PublishedMax,
		input.PeriodMin, input.PeriodMax,
		input.Tags)
	results := this.documents.Search(criteria)
	return jsonResult(results)
}
func (this *Search) Tags(input *inputs.SearchTag) detour.Renderer {
	criteria := projections.NewTagCriteria(input.Text, input.Tags)
	results := this.tags.Search(criteria)
	return jsonResult(results)
}

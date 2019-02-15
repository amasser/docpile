package http

import (
	"github.com/joliver/docpile/app/http/inputs"
	"github.com/joliver/docpile/app/projections"
	"github.com/smartystreets/detour"
)

type Search struct {
	documents *projections.AllDocuments
	tags      *projections.MatchingTags
}

func NewSearch(documents *projections.AllDocuments, tags *projections.MatchingTags) *Search {
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
	results := this.tags.Search(input.Text, input.Tags)
	return jsonResult(results)
}

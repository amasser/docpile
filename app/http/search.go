package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/app/projections"
	"github.com/smartystreets/detour"
)

type Search struct {
	allDocuments  *projections.AllDocuments
	tagProjection *projections.TagProjection
}

func NewSearch(allDocuments *projections.AllDocuments, tagProjection *projections.TagProjection) *Search {
	return &Search{allDocuments: allDocuments, tagProjection: tagProjection}
}

func (this *Search) Documents(input *inputs.SearchDocument) detour.Renderer {
	criteria := projections.NewDocumentCriteria(
		input.PublishedMin, input.PublishedMax,
		input.PeriodMin, input.PeriodMax,
		input.Tags)
	results := this.allDocuments.Search(criteria)
	return jsonResult(results)
}
func (this *Search) Tags(input *inputs.SearchTag) detour.Renderer {
	criteria := projections.NewTagCriteria(input.Text, input.Tags)
	results := this.tagProjection.Search(criteria)
	return jsonResult(results)
}

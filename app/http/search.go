package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/app/projections"
	"github.com/smartystreets/detour"
)

type Search struct {
	projection *projections.AllDocuments
}

func NewSearch(projection *projections.AllDocuments) *Search {
	return &Search{projection: projection}
}

func (this *Search) Documents(input *inputs.SearchDocument) detour.Renderer {
	return jsonResult(this.projection.Search(projections.NewDocumentCriteria(
		input.PublishedMin, input.PublishedMax,
		input.PeriodMin, input.PeriodMax,
		input.Tags)))
}
func (this *Search) Tags(input *inputs.SearchTag) detour.Renderer {
	return nil
}

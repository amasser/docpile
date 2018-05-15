package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"bitbucket.org/jonathanoliver/docpile/app/projections"
	"github.com/smartystreets/detour"
)

type Search struct {
	projector *projections.Projector
}

func NewSearch(projector *projections.Projector) *Search {
	return &Search{projector: projector}
}

func (this *Search) Documents(input *inputs.SearchDocument) detour.Renderer {
	return jsonResult(this.projector.SearchDocuments(projections.NewDocumentSearch(
		input.PublishedMin, input.PublishedMax,
		input.PeriodMin, input.PeriodMax,
		input.Tags)))
}
func (this *Search) Tags(input *inputs.SearchTag) detour.Renderer {
	return nil
}

package http

import (
	"bitbucket.org/jonathanoliver/docpile/app/http/inputs"
	"github.com/smartystreets/detour"
)

type Search struct {
}

func NewSearch(_ interface{}) *Search {
	return &Search{ /* TODO */ }
}

func (this *Search) Documents(input *inputs.SearchDocument) detour.Renderer {
	return nil
}

func (this *Search) Tags(input *inputs.SearchTag) detour.Renderer {
	return nil
}

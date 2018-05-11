package domain

import (
	"bitbucket.org/jonathanoliver/docpile/infrastructure"
)

type Handler struct {
	aggregate  infrastructure.Aggregate
	applicator infrastructure.Applicator
}

func NewHandler(aggregate infrastructure.Aggregate, applicator infrastructure.Applicator) *Handler {
	return &Handler{aggregate: aggregate, applicator: applicator}
}

func (this *Handler) Handle(message interface{}) infrastructure.Result {
	result := this.aggregate.Handle(message)
	if result.Error == nil {
		this.applicator.Apply(this.aggregate.Consume())
	}

	return result
}

package handlers

import "bitbucket.org/jonathanoliver/docpile/infrastructure"

type DomainHandler struct {
	aggregate  infrastructure.Aggregate
	applicator infrastructure.Applicator
}

func NewDomainHandler(aggregate infrastructure.Aggregate, applicator infrastructure.Applicator) *DomainHandler {
	return &DomainHandler{aggregate: aggregate, applicator: applicator}
}

func (this *DomainHandler) Handle(message interface{}) infrastructure.Result {
	result := this.aggregate.Handle(message)
	this.applicator.Apply(this.aggregate.Consume())
	return result
}

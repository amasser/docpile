package handlers

import "bitbucket.org/jonathanoliver/docpile/generic"

type DomainHandler struct {
	aggregate  generic.Aggregate
	applicator generic.Applicator
}

func NewDomainHandler(aggregate generic.Aggregate, applicator generic.Applicator) *DomainHandler {
	return &DomainHandler{aggregate: aggregate, applicator: applicator}
}

func (this *DomainHandler) Handle(message interface{}) generic.Result {
	result := this.aggregate.Handle(message)
	this.applicator.Apply(this.aggregate.Consume())
	return result
}

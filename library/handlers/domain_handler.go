package handlers

import "bitbucket.org/jonathanoliver/docpile/library"

type DomainHandler struct {
	aggregate  library.Aggregate
	applicator library.Applicator
}

func NewDomainHandler(aggregate library.Aggregate, applicator library.Applicator) *DomainHandler {
	return &DomainHandler{aggregate: aggregate, applicator: applicator}
}

func (this *DomainHandler) Handle(message interface{}) library.Result {
	result := this.aggregate.Handle(message)
	this.applicator.Apply(this.aggregate.Consume())
	return result
}

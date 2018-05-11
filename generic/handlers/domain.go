package handlers

import "bitbucket.org/jonathanoliver/docpile/generic/applicators"

type Domain struct {
	aggregate  Aggregate
	applicator applicators.Applicator
}

func NewDomain(aggregate Aggregate, applicator applicators.Applicator) *Domain {
	return &Domain{aggregate: aggregate, applicator: applicator}
}

func (this *Domain) Handle(message interface{}) Result {
	result := this.aggregate.Handle(message)
	this.applicator.Apply(this.aggregate.Consume())
	return result
}

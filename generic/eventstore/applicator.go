package eventstore

import "bitbucket.org/jonathanoliver/docpile/generic"

type Applicator struct {
	inner generic.Applicator
	store EventStore
}

func NewApplicator(inner generic.Applicator, store EventStore) *Applicator {
	return &Applicator{inner: inner, store: store}
}

func (this *Applicator) Apply(messages []interface{}) {
	if err := this.store.Store(messages); err != nil {
		panic(err)
	}

	this.inner.Apply(messages)
}

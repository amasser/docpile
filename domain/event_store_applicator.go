package domain

type EventStoreApplicator struct {
	inner Applicator
	store EventStore
}

func NewEventStoreApplicator(inner Applicator, store EventStore) *EventStoreApplicator {
	return &EventStoreApplicator{inner: inner, store: store}
}

func (this *EventStoreApplicator) Apply(messages []interface{}) {
	if err := this.store.Store(messages); err != nil {
		panic(err)
	}

	this.inner.Apply(messages)
}

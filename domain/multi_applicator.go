package domain

type MultiApplicator struct {
	inner []Applicator
}

func NewMultiApplicator(inner ...Applicator) *MultiApplicator {
	return &MultiApplicator{inner: inner}
}

func (this *MultiApplicator) Apply(messages []interface{}) {
	for _, inner := range this.inner {
		inner.Apply(messages)
	}
}

package applicators

type Fanout struct {
	inner []Applicator
}

func NewFanout(inner ...Applicator) *Fanout {
	return &Fanout{inner: inner}
}

func (this *Fanout) Apply(messages []interface{}) {
	for _, inner := range this.inner {
		inner.Apply(messages)
	}
}

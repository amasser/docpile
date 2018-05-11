package projections

type Projector struct {
	tags *AllTags
}

func NewProjector() *Projector {
	return &Projector{
		tags: NewAllTags(),
	}
}

func (this *Projector) Apply(messages []interface{}) {
	for _, message := range messages {
		this.apply(message)
	}
}
func (this *Projector) apply(message interface{}) {
	this.tags.Transform(message)
}

func (this *Projector) ListTags() interface{} { return this.tags }

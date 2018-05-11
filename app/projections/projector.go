package projections

type Projector struct {
	allTags *AllTags
}

func NewProjector() *Projector {
	return &Projector{
		allTags: NewAllTags(),
	}
}

func (this *Projector) Apply(messages []interface{}) {
	for _, message := range messages {
		this.apply(message)
	}
}
func (this *Projector) apply(message interface{}) {
	this.allTags.Transform(message)
}

func (this *Projector) ListTags() interface{} { return this.allTags.List() }

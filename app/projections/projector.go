package projections

type Projector struct {
	AllTags      *AllTags
	AllDocuments *AllDocuments
	TagSearch    *TagSearch
}

func NewProjector() *Projector {
	return &Projector{
		AllTags:      NewAllTags(),
		AllDocuments: NewAllDocuments(),
		TagSearch:    NewTagSearch(),
	}
}

///////////////////////////////////////////

func (this *Projector) Apply(messages []interface{}) {
	for _, message := range messages {
		this.apply(message)
	}
}
func (this *Projector) apply(message interface{}) {
	this.AllTags.Transform(message)
	this.AllDocuments.Transform(message)
	this.TagSearch.Transform(message)
}

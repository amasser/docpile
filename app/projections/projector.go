package projections

type Projector struct {
	allTags      *AllTags
	allDocuments *AllDocuments
}

func NewProjector() *Projector {
	return &Projector{
		allTags:      NewAllTags(),
		allDocuments: NewAllDocuments(),
	}
}

func (this *Projector) Apply(messages []interface{}) {
	for _, message := range messages {
		this.apply(message)
	}
}
func (this *Projector) apply(message interface{}) {
	this.allTags.Transform(message)
	this.allDocuments.Transform(message)
}

func (this *Projector) ListTags() interface{}                       { return this.allTags.List() }
func (this *Projector) LoadTag(id uint64) (interface{}, error)      { return this.allTags.Load(id) }
func (this *Projector) ListDocuments() interface{}                  { return this.allDocuments.List() }
func (this *Projector) LoadDocument(id uint64) (interface{}, error) { return this.allDocuments.Load(id) }

func (this *Projector) SearchDocuments(search DocumentSpecification) interface{} {
	var documents []Document

	for _, document := range this.allDocuments.List() {
		if search.IsSatisfiedBy(document) {
			documents = append(documents, document)
		}
	}

	return documents
}

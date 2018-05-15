package projections

type TagCriteria struct {
	text string
	tags []uint64
}

func NewTagCriteria(text string, tags []uint64) TagCriteria {
	return TagCriteria{text: text, tags: tags}
}

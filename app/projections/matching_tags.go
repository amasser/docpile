package projections

type MatchingTags struct {
	documents *AllDocuments
	tags      *AllTags
}

func NewMatchingTags(documents *AllDocuments, tags *AllTags) *MatchingTags {
	return &MatchingTags{documents: documents, tags: tags}
}

func (this *MatchingTags) Search(text string, tags []uint64) []MatchingTag {
	search := NewTagSearch(this.documents.items, this.tags.index, this.tags.items)
	return search.Search(text, tags)
}

type MatchingTag struct {
	TagID   uint64 `json:"tag_id"`
	TagText string `json:"text"`
	Synonym bool   `json:"synonym"`
	Indexes []int  `json:"indexes"`
}

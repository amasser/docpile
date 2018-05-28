package projections

import "bitbucket.org/jonathanoliver/docpile/app/events"

type MatchingTags struct {
}

func NewMatchingTags() *MatchingTags {
	return &MatchingTags{}
}

func (this *MatchingTags) Transform(message interface{}) {
	switch message := message.(type) {
	case events.TagAdded:
		this.tagAdded(message)
	case events.TagRemoved:
		this.tagRemoved(message)
	case events.TagRenamed:
		this.tagRenamed(message)
	case events.TagSynonymDefined:
		this.synonymDefined(message)
	case events.TagSynonymRemoved:
		this.synonymRemoved(message)
	case events.DocumentDefined:
		this.documentDefined(message)
	case events.DocumentRemoved:
		this.documentRemoved(message)
	}
}
func (this *MatchingTags) tagAdded(message events.TagAdded) {
}
func (this *MatchingTags) tagRemoved(message events.TagRemoved) {
}
func (this *MatchingTags) tagRenamed(message events.TagRenamed) {
}
func (this *MatchingTags) synonymDefined(message events.TagSynonymDefined) {
}
func (this *MatchingTags) synonymRemoved(message events.TagSynonymRemoved) {
}
func (this *MatchingTags) documentDefined(message events.DocumentDefined) {
}
func (this *MatchingTags) documentRemoved(message events.DocumentRemoved) {
}

func (this *MatchingTags) Search(text string, tags []uint64) []MatchingTag {
	search := NewTagSearch(text, tags)
	return search.Search()
}

type MatchingTag struct {
	TagID   uint64 `json:"tag_id"`
	TagText string `json:"text"`
	Synonym bool   `json:"synonym"`
	Indexes []int  `json:"indexes"`
}

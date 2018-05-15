package projections

import "bitbucket.org/jonathanoliver/docpile/app/events"

type TagSearch struct {
}

func NewTagSearch() *TagSearch {
	return &TagSearch{}
}

func (this *TagSearch) Transform(message interface{}) {
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
func (this *TagSearch) tagAdded(message events.TagAdded) {
}
func (this *TagSearch) tagRemoved(message events.TagRemoved) {
}
func (this *TagSearch) tagRenamed(message events.TagRenamed) {
}
func (this *TagSearch) synonymDefined(message events.TagSynonymDefined) {
}
func (this *TagSearch) synonymRemoved(message events.TagSynonymRemoved) {
}
func (this *TagSearch) documentDefined(message events.DocumentDefined) {
}
func (this *TagSearch) documentRemoved(message events.DocumentRemoved) {
}

func (this *TagSearch) Search(criteria TagCriteria) (matching []MatchingTag) {
	return nil
}

type MatchingTag struct {
	TagID   uint64 `json:"tag_id"`
	TagText string `json:"text"`
	Synonym bool   `json:"synonym"`
}

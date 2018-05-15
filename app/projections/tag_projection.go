package projections

import "bitbucket.org/jonathanoliver/docpile/app/events"

type TagProjection struct {
}

func NewTagProjection() *TagProjection {
	return &TagProjection{}
}

func (this *TagProjection) Transform(message interface{}) {
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
func (this *TagProjection) tagAdded(message events.TagAdded) {
}
func (this *TagProjection) tagRemoved(message events.TagRemoved) {
}
func (this *TagProjection) tagRenamed(message events.TagRenamed) {
}
func (this *TagProjection) synonymDefined(message events.TagSynonymDefined) {
}
func (this *TagProjection) synonymRemoved(message events.TagSynonymRemoved) {
}
func (this *TagProjection) documentDefined(message events.DocumentDefined) {
}
func (this *TagProjection) documentRemoved(message events.DocumentRemoved) {
}

func (this *TagProjection) Search(criteria TagCriteria) (matching []MatchingTag) {
	return nil
}

type MatchingTag struct {
	TagID   uint64 `json:"tag_id"`
	TagText string `json:"text"`
	Synonym bool   `json:"synonym"`
}

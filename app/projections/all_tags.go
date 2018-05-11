package projections

import (
	"time"

	"bitbucket.org/jonathanoliver/docpile/app/events"
)

type AllTags struct {
	index map[uint64]int
	tags  []Tag
}

func NewAllTags() *AllTags {
	return &AllTags{index: map[uint64]int{}}
}

func (this *AllTags) loadTag(tagID uint64) *Tag {
	if index, contains := this.index[tagID]; contains {
		return &this.tags[index]
	} else {
		return &Tag{Synonyms: map[string]time.Time{}}
	}
}

func (this *AllTags) Transform(message interface{}) {
	switch message := message.(type) {
	case events.TagAdded:
		this.tagAdded(message)
	case events.TagRenamed:
		this.tagRenamed(message)
	case events.TagSynonymDefined:
		this.synonymDefined(message)
	case events.TagSynonymRemoved:
		this.synonymRemoved(message)
	}
}

func (this *AllTags) tagAdded(message events.TagAdded) {
	if _, contains := this.index[message.TagID]; !contains {
		this.tags = append(this.tags, newTag(message))
		this.index[message.TagID] = len(this.tags) - 1
	}
}
func (this *AllTags) tagRenamed(message events.TagRenamed) {
	this.loadTag(message.TagID).TagName = message.NewName
}
func (this *AllTags) synonymDefined(message events.TagSynonymDefined) {
	this.loadTag(message.TagID).Synonyms[message.TagName] = message.Timestamp
}
func (this *AllTags) synonymRemoved(message events.TagSynonymRemoved) {
	delete(this.loadTag(message.TagID).Synonyms, message.TagName)
}

func (this *AllTags) List() []Tag { return this.tags }

//////////////////////////////////////////////////////////////

type Tag struct {
	TagID     uint64               `json:"tag_id"`
	Timestamp time.Time            `json:"timestamp"`
	TagName   string               `json:"tag_name"`
	Synonyms  map[string]time.Time `json:"synonyms,omitempty"`
}

func newTag(message events.TagAdded) Tag {
	return Tag{
		TagID:     message.TagID,
		Timestamp: message.Timestamp,
		TagName:   message.TagName,
		Synonyms:  map[string]time.Time{},
	}
}

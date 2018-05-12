package projections

import (
	"errors"
	"time"

	"bitbucket.org/jonathanoliver/docpile/app/events"
)

type AllTags struct {
	index map[uint64]int
	items []Tag
}

func NewAllTags() *AllTags {
	return &AllTags{index: map[uint64]int{}}
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
		this.index[message.TagID] = len(this.items)
		this.items = append(this.items, newTag(message))
	}
}
func (this *AllTags) tagRenamed(message events.TagRenamed) {
	this.load(message.TagID).TagName = message.NewName
}
func (this *AllTags) synonymDefined(message events.TagSynonymDefined) {
	this.load(message.TagID).Synonyms[message.TagName] = message.Timestamp
}
func (this *AllTags) synonymRemoved(message events.TagSynonymRemoved) {
	delete(this.load(message.TagID).Synonyms, message.TagName)
}

func (this *AllTags) load(id uint64) *Tag {
	if index, contains := this.index[id]; contains {
		return &this.items[index]
	} else {
		return &Tag{Synonyms: map[string]time.Time{}}
	}
}

func (this *AllTags) List() []Tag { return this.items }
func (this *AllTags) Load(id uint64) (Tag, error) {
	if index, contains := this.index[id]; contains {
		return this.items[index], nil
	} else {
		return Tag{}, TagNotFoundError
	}
}

var (
	TagNotFoundError = errors.New("tag not found")
)

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

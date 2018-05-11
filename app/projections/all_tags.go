package projections

import (
	"time"

	"bitbucket.org/jonathanoliver/docpile/app/events"
)

type AllTags struct {
	Tags []*Tag `json:"tags,omitempty"`
}

func NewAllTags() *AllTags {
	return &AllTags{}
}

func (this *AllTags) Apply(messages []interface{}) {
	for _, message := range messages {
		this.apply(message)
	}
}
func (this *AllTags) apply(message interface{}) {
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
	for _, tag := range this.Tags {
		if tag.TagID == message.TagID {
			return
		}
	}

	this.Tags = append(this.Tags, newTag(message))
}
func (this *AllTags) tagRenamed(message events.TagRenamed) {
	for _, tag := range this.Tags {
		if tag.TagID != message.TagID {
			continue
		}

		tag.TagName = message.NewName
		break
	}
}
func (this *AllTags) synonymDefined(message events.TagSynonymDefined) {
	for _, tag := range this.Tags {
		if tag.TagID != message.TagID {
			continue
		}

		tag.Synonyms[message.TagName] = message.Timestamp
		break
	}
}
func (this *AllTags) synonymRemoved(message events.TagSynonymRemoved) {
	for _, tag := range this.Tags {
		if tag.TagID != message.TagID {
			continue
		}

		delete(tag.Synonyms, message.TagName)
		break
	}
}

type Tag struct {
	TagID     uint64               `json:"tag_id"`
	Timestamp time.Time            `json:"timestamp"`
	TagName   string               `json:"tag_name"`
	Synonyms  map[string]time.Time `json:"synonyms,omitempty"`
}

func newTag(message events.TagAdded) *Tag {
	return &Tag{
		TagID:     message.TagID,
		Timestamp: message.Timestamp,
		TagName:   message.TagName,
		Synonyms:  map[string]time.Time{},
	}
}

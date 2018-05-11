package events

import (
	"log"
	"reflect"
	"time"
)

type SHA256Hash [32]byte

type TagAdded struct {
	TagID     uint64    `json:"tag_id"`
	Timestamp time.Time `json:"timestamp"`
	TagName   string    `json:"tag_name"`
}
type TagRenamed struct {
	TagID     uint64    `json:"tag_id"`
	Timestamp time.Time `json:"timestamp"`
	OldName   string    `json:"old_name"`
	NewName   string    `json:"new_name"`
}
type TagSynonymDefined struct {
	TagID     uint64    `json:"tag_id"`
	Timestamp time.Time `json:"timestamp"`
	TagName   string    `json:"tag_name"`
}
type TagSynonymRemoved struct {
	TagID     uint64    `json:"tag_id"`
	Timestamp time.Time `json:"timestamp"`
	TagName   string    `json:"tag_name"`
}

type ManagedAssetImported struct {
	AssetID   uint64     `json:"asset_id"`
	Timestamp time.Time  `json:"timestamp"`
	Hash      SHA256Hash `json:"sha256"`
	MIMEType  string     `json:"mime_type"`
	Name      string     `json:"name"`
	Key       string     `json:"key"`
}

// NOTE: this could simply be for resources which are unmanaged rather than cloud
type CloudAssetImported struct {
	AssetID   uint64    `json:"asset_id"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name,omitempty"`
	Provider  string    `json:"provider"`
	Resource  string    `json:"resource"`
}

type DocumentDefined struct {
	DocumentID  uint64     `json:"document_id"`
	Timestamp   time.Time  `json:"timestamp"`
	AssetID     uint64     `json:"asset_id"`
	AssetOffset uint64     `json:"asset_offset,omitempty"`
	Published   *time.Time `json:"published,omitempty"`
	PeriodBegin *time.Time `json:"begin,omitempty"`
	PeriodEnd   *time.Time `json:"end,omitempty"`
	Tags        []uint64   `json:"tags,omitempty"`
	Documents   []uint64   `json:"documents,omitempty"`
	Description string     `json:"description,omitempty"`
}

///////////////////////////////////////////////////////////

func init() {
	registerType("tag-added", TagAdded{})
	registerType("tag-removed", TagRenamed{})
	registerType("tag-synonym-defined", TagSynonymDefined{})
	registerType("tag-synonym-removed", TagSynonymRemoved{})
	registerType("managed-asset-imported", ManagedAssetImported{})
	registerType("cloud-asset-imported", CloudAssetImported{})
	registerType("document-defined", DocumentDefined{})
}
func registerType(name string, instance interface{}) {
	instanceType := reflect.TypeOf(instance)
	registry[name] = instanceType
}
func InstanceRegistry(name string) reflect.Type {
	if registeredType, contains := registry[name]; contains {
		return registeredType
	} else {
		log.Fatalf("Unknown message type: [%s]\n", name)
		return nil
	}
}

var registry = map[string]reflect.Type{}

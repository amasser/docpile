package events

import (
	"time"
)

type SHA256Hash [32]byte

type TagAdded struct {
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
}
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
	Published   *time.Time `json:"published,omitempty"`
	PeriodBegin *time.Time `json:"begin,omitempty"`
	PeriodEnd   *time.Time `json:"end,omitempty"`
	Tags        []uint64   `json:"tags,omitempty"`
	Associated  []uint64   `json:"Associated,omitempty"`
	Description string     `json:"description,omitempty"`
}

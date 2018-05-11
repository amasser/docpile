package domain

import "time"

type DocumentDefinition struct {
	AssetID     uint64
	AssetOffset uint64
	Published   *time.Time
	PeriodBegin *time.Time
	PeriodEnd   *time.Time
	Tags        []uint64
	Documents   []uint64
	Description string
}

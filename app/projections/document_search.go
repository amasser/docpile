package projections

import "time"

type DocumentSearch struct {
	publishedMin *time.Time
	publishedMax *time.Time
	periodMin    *time.Time
	periodMax    *time.Time
	tags         []uint64
}

func NewDocumentSearch(
	publishedMin *time.Time,
	publishedMax *time.Time,
	periodMin *time.Time,
	periodMax *time.Time,
	tags []uint64,
) *DocumentSearch {
	return &DocumentSearch{
		publishedMin: publishedMin,
		publishedMax: publishedMax,
		periodMin:    periodMin,
		periodMax:    periodMax,
		tags:         tags,
	}
}

func (this *DocumentSearch) IsSatisfiedBy(document Document) bool {
	return true
}

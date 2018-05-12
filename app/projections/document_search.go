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
	// TODO: period min/max
	return this.withinPublishedLimits(document.Published) &&
		documentContainsAllSearchTags(this.tags, document.Tags)
}

func (this *DocumentSearch) withinPublishedLimits(published *time.Time) bool {
	if published == nil {
		return true
	}

	if this.publishedMin != nil && this.publishedMin.After(*published) {
		return false
	}

	if this.publishedMax != nil && this.publishedMax.Before(*published) {
		return false
	}

	return true
}

func documentContainsAllSearchTags(searchTags, documentTags []uint64) bool {
	for _, searchTag := range searchTags {
		if !documentTagsContainSearchTag(searchTag, documentTags) {
			return false
		}
	}

	return true
}
func documentTagsContainSearchTag(searchTag uint64, documentTags []uint64) bool {
	for _, documentTag := range documentTags {
		if documentTag == searchTag {
			return true
		}
	}

	return false
}

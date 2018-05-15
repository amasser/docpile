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
) DocumentSearch {
	return DocumentSearch{
		publishedMin: publishedMin,
		publishedMax: publishedMax,
		periodMin:    periodMin,
		periodMax:    periodMax,
		tags:         tags,
	}
}

func (this *DocumentSearch) IsSatisfiedBy(document Document) bool {
	return this.withinPublishedLimits(document.Published) &&
		this.withinPeriodLimits(document.PeriodMin, document.PeriodMax) &&
		this.containsAllSearchTags(document.Tags)
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

func (this *DocumentSearch) withinPeriodLimits(min, max *time.Time) bool {
	if min == nil && max == nil {
		return true
	}

	if min != nil && this.periodMin != nil && this.periodMin.After(*min) {
		return false
	}

	if max != nil && this.periodMin != nil && this.periodMax.Before(*max) {
		return false
	}

	return true
}

func (this *DocumentSearch) containsAllSearchTags(documentTags []uint64) bool {
	for _, searchTag := range this.tags {
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

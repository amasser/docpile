package inputs

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/smartystreets/detour"
)

type DocumentSearch struct {
	PublishedMin *time.Time `json:"published_min"`
	PublishedMax *time.Time `json:"published_max"`
	PeriodMin    *time.Time `json:"period_min"`
	PeriodMax    *time.Time `json:"period_max"`
	Tags         []uint64   `json:"tags"`
}

func (this *DocumentSearch) Bind(request *http.Request) error {
	return json.NewDecoder(request.Body).Decode(this)
}

func (this *DocumentSearch) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(duplicateTagError, containsDuplicate(this.Tags))
	errors = errors.AppendIf(noDocumentSearchCriteria, this.emptySearch())
	return errors
}
func (this *DocumentSearch) emptySearch() bool {
	return this.PublishedMin == nil && this.PublishedMax == nil &&
		this.PeriodMin == nil && this.PeriodMax == nil &&
		len(this.Tags) == 0
}

var noDocumentSearchCriteria = detour.CompoundInputError(
	"At least one document search criteria must be provided.",
	jsonPeriodMinField, jsonPeriodMaxField, jsonPublishedMinField, jsonPublishedMaxField)

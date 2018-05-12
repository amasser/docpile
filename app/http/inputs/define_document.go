package inputs

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/smartystreets/detour"
)

type DefineDocument struct {
	AssetID     uint64     `json:"asset_id"`
	AssetOffset uint64     `json:"asset_offset"`
	Published   *time.Time `json:"published"`
	PeriodMin   *time.Time `json:"period_min"`
	PeriodMax   *time.Time `json:"period_max"`
	Tags        []uint64   `json:"tags"`
	Documents   []uint64   `json:"documents"`
	Description string     `json:"description"`
}

func (this *DefineDocument) Bind(request *http.Request) error {
	return json.NewDecoder(request.Body).Decode(this)
}

func (this *DefineDocument) Sanitize() {
	this.Description = strings.TrimSpace(this.Description)
}

func (this *DefineDocument) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(minMustComeFirstError,
		this.PeriodMin != nil && this.PeriodMax != nil && this.PeriodMin.After(*this.PeriodMax))
	errors = errors.AppendIf(maxMustHaveAMinError, this.PeriodMax != nil && this.PeriodMin == nil)
	errors = errors.AppendIf(duplicateTagError, containsDuplicate(this.Tags))
	errors = errors.AppendIf(duplicateDocumentError, containsDuplicate(this.Documents))
	return errors
}

func containsDuplicate(values []uint64) bool {
	unique := make(map[uint64]struct{}, len(values))
	for _, value := range values {
		if _, contains := unique[value]; contains {
			return true
		} else {
			unique[value] = struct{}{}
		}
	}

	return false
}

var (
	minMustComeFirstError      = fieldError("The min/begin date must come on or before the max/end date.", jsonPeriodMinField, jsonPeriodMaxField)
	maxMustHaveAMinError       = fieldError("When an max/end date is specified it must have a min/begin date.", jsonPeriodMinField, jsonPeriodMaxField)
	duplicateTagError          = fieldError("The tag IDs provided must be unique.", jsonTagsField)
	duplicateDocumentError     = fieldError("The document IDs provided must be unique.", jsonDocumentsField)
	AssetDoesNotExistResult    = notFoundResult("The asset ID supplied could not be found.", jsonAssetIDField)
	TagDoesNotExistResult      = notFoundResult("One or more tag supplied could not be found.", jsonTagsField)
	DocumentDoesNotExistResult = notFoundResult("One or more documents supplied could not be found.", jsonDocumentsField)
)

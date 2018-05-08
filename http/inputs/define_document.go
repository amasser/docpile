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
	PeriodBegin *time.Time `json:"period_begin"`
	PeriodEnd   *time.Time `json:"period_end"`
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
	errors = errors.AppendIf(beforeMustComeFirstError,
		this.PeriodBegin != nil && this.PeriodEnd != nil && this.PeriodBegin.After(*this.PeriodEnd))
	errors = errors.AppendIf(endMustHaveABeforeError, this.PeriodEnd != nil && this.PeriodBegin == nil)
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
	beforeMustComeFirstError   = fieldError("The begin date must come before the end date.", jsonPeriodBeginField, jsonPeriodEndField)
	endMustHaveABeforeError    = fieldError("When an end date is specified it must have a begin date.", jsonPeriodBeginField, jsonPeriodEndField)
	duplicateTagError          = fieldError("The tag IDs provided must be unique.", jsonTagsField)
	duplicateDocumentError     = fieldError("The document IDs provided must be unique.", jsonDocumentsField)
	AssetDoesNotExistResult    = notFoundResult("The asset ID supplied could not be found.", jsonAssetIDField)
	TagDoesNotExistResult      = notFoundResult("One or more tag supplied could not be found.", jsonTagsField)
	DocumentDoesNotExistResult = notFoundResult("One or more documents supplied could not be found.", jsonDocumentsField)
)

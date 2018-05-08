package inputs

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/smartystreets/detour"
)

type DefineDocument struct {
	FileID      uint64     `json:"file_id"`
	FileIndex   uint64     `json:"file_index"`
	Date        *time.Time `json:"date"`
	Begin       *time.Time `json:"begin"`
	End         *time.Time `json:"end"`
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
		this.Begin != nil && this.End != nil && this.Begin.After(*this.End))
	errors = errors.AppendIf(endMustHaveABeforeError, this.End != nil && this.Begin == nil)

	return errors
}

func readNumbers(query url.Values, key string) (numbers []uint64, err error) {
	for _, item := range query[key] {
		if parsed, err := strconv.ParseUint(item, 10, 64); err != nil {
			return nil, err
		} else {
			numbers = append(numbers, parsed)
		}
	}

	return numbers, err
}

var (
	beforeMustComeFirstError   = fieldError("The begin date must come before the end date.", beginField, endField)
	endMustHaveABeforeError    = fieldError("When an end date is specified it must have a begin date.", beginField, endField)
	AssetDoesNotExistResult    = notFoundResult("The asset ID supplied could not be found.", jsonAssetIDField)
	TagDoesNotExistResult      = notFoundResult("One or more tag supplied could not be found.", jsonTagsField)
	DocumentDoesNotExistResult = notFoundResult("One or more documents supplied could not be found.", jsonDocumentsField)
)

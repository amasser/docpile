package inputs

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/smartystreets/detour"
)

type SearchTag struct {
	Text string   `json:"text"`
	Tags []uint64 `json:"tags"`
}

func (this *SearchTag) Bind(request *http.Request) error {
	return json.NewDecoder(request.Body).Decode(this)
}

func (this *SearchTag) Sanitize() {
	this.Text = strings.TrimSpace(this.Text)
}

func (this *SearchTag) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(noTagSearchCriteria, len(this.Text) == 0)
	errors = errors.AppendIf(duplicateTagError, containsDuplicate(this.Tags))
	return errors
}

var noTagSearchCriteria = detour.SimpleInputError("No query text provided.", jsonTextField)

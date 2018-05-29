package inputs

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/smartystreets/detour"
)

type TagInput struct {
	ID   uint64 `json:"-"`
	Name string `json:"name"`
}

func (this *TagInput) Bind(request *http.Request) error {
	this.ID = idFromURLPath(request)
	return json.NewDecoder(request.Body).Decode(this)
}

func (this *TagInput) Sanitize() {
	this.Name = strings.TrimSpace(this.Name)
}

func (this *TagInput) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(missingTagIDError, this.ID == 0)
	errors = errors.AppendIf(missingTagNameError, len(this.Name) == 0)
	return errors
}

var (
	missingTagIDError     = fieldError("A tag ID is required.", jsonTagIDField)
	TagNotFoundResult     = notFoundResult("The tag ID was not found.", jsonTagIDField)
	SynonymNotFoundResult = notFoundResult("The tag synonym was not found.", jsonNameField)
)

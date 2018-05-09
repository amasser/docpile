package inputs

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/smartystreets/detour"
)

type RenameTag struct {
	ID   uint64 `json:"tag_id"`
	Name string `json:"name"`
}

func (this *RenameTag) Bind(request *http.Request) error {
	return json.NewDecoder(request.Body).Decode(this)
}

func (this *RenameTag) Sanitize() {
	this.Name = strings.TrimSpace(this.Name)
}

func (this *RenameTag) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(missingTagIDError, this.ID == 0)
	errors = errors.AppendIf(missingTagNameError, len(this.Name) == 0)
	return errors
}

var (
	missingTagIDError = fieldError("A tag ID is required.", jsonTagIDField)
	TagNotFoundResult = notFoundResult("", jsonTagIDField)
)

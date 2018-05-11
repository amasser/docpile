package inputs

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/smartystreets/detour"
)

type AddTag struct {
	Name string `json:"name"`
}

func (this *AddTag) Bind(request *http.Request) error {
	return json.NewDecoder(request.Body).Decode(this)
}

func (this *AddTag) Sanitize() {
	this.Name = strings.TrimSpace(this.Name)
}

func (this *AddTag) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(missingTagNameError, len(this.Name) == 0)
	return errors
}

var (
	missingTagNameError = fieldError("A tag name is required.", jsonNameField)
	DuplicateTagResult  = conflictResult("The tag provided already exists.", jsonNameField)
)

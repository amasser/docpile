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
	errors = errors.AppendIf(missingTagError, len(this.Name) == 0)
	return errors
}

var (
	missingTagError    = detour.SimpleInputError("A tag name is required.", nameField)
	DuplicateTagResult = detour.ErrorResult{
		StatusCode: http.StatusConflict,
		Error1:     detour.SimpleInputError("The tag provided already exists.", nameField),
	}
)

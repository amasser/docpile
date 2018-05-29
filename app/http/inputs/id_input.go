package inputs

import (
	"net/http"

	"github.com/smartystreets/detour"
)

type IDInput struct {
	ID uint64
}

func (this *IDInput) Bind(request *http.Request) error {
	this.ID = idFromURLPath(request)
	return nil
}

func (this *IDInput) Validate() error {
	var errors detour.Errors
	errors = errors.AppendIf(idNotFoundError, this.ID == 0)
	return errors
}

var (
	idNotFoundError  = detour.SimpleInputError("An ID must be specified.", urlIDField)
	IDNotFoundResult = detour.ErrorResult{
		StatusCode: http.StatusNotFound,
		Error1:     detour.SimpleInputError("The ID provided was not found.", urlIDField),
	}
)

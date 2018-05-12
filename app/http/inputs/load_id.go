package inputs

import (
	"net/http"
	"path"
	"strconv"

	"github.com/smartystreets/detour"
)

type LoadID struct {
	ID uint64
}

func (this *LoadID) Bind(request *http.Request) error {
	this.ID, _ = strconv.ParseUint(path.Base(request.URL.Path), 10, 64)
	return nil
}

func (this *LoadID) Validate() error {
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

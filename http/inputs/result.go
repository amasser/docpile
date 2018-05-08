package inputs

import (
	"net/http"

	"github.com/smartystreets/detour"
)

func conflictResult(message string, fields ...string) detour.Renderer {
	return detour.ErrorResult{
		StatusCode: http.StatusConflict,
		Error1:     fieldError(message, fields...),
	}
}
func notFoundResult(message string, fields ...string) detour.Renderer {
	return detour.ErrorResult{
		StatusCode: http.StatusNotFound,
		Error1:     fieldError(message, fields...),
	}
}
func fieldError(message string, fields ...string) error {
	if len(fields) == 0 {
		return detour.SimpleInputError(message, "")
	} else if len(fields) == 1 {
		return detour.SimpleInputError(message, fields[0])
	} else {
		return detour.CompoundInputError(message, fields...)
	}
}

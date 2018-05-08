package http

import (
	"net/http"

	"github.com/smartystreets/detour"
)

var UnknownErrorResult = detour.StatusCodeResult{
	StatusCode: http.StatusInternalServerError,
	Message:    "Unhandled internal error",
}

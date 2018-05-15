package http

import (
	"net/http"

	"github.com/smartystreets/detour"
)

var UnknownErrorResult = detour.StatusCodeResult{
	StatusCode: http.StatusInternalServerError,
	Message:    "Unhandled internal error",
}

func newEntityResult(id uint64) detour.Renderer {
	return detour.JSONResult{Content: id, StatusCode: http.StatusCreated}
}

func jsonResult(value interface{}) detour.Renderer {
	return detour.JSONResult{Content: value}
}

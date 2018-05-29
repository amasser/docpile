package inputs

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func fromURLPath(request *http.Request, name string) string {
	context := request.Context()
	parameters := httprouter.ParamsFromContext(context)
	return parameters.ByName(name)
}

func idFromURLPath(request *http.Request) uint64 {
	id, _ := strconv.ParseUint(fromURLPath(request, urlIDField), 10, 64)
	return id
}

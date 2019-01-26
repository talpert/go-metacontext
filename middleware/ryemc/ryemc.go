package ryemc

import (
	"net/http"

	"github.com/InVisionApp/rye"
	"github.com/talpert/go-metacontext"
)

// Rye middleware to parse the request
func ParseMetaContextMiddleware(rw http.ResponseWriter, r *http.Request) *rye.Response {
	ctx, err := metacontext.ParseRequest(r)
	if err != nil {
		return &rye.Response{
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	return &rye.Response{
		Context: ctx,
	}
}

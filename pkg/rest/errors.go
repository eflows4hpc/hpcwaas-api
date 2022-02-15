package rest

import (
	"fmt"
	"net/http"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/gin-gonic/gin"
)

func writeError(gc *gin.Context, err *api.Error) {
	gc.JSON(err.Status, api.Errors{Errors: []*api.Error{err}})
	gc.AbortWithStatus(err.Status)
}

var (
	errNotFound  = &api.Error{ID: "not_found", Status: 404, Title: "Not Found", Detail: "Requested content not found."}
	errForbidden = &api.Error{ID: "forbidden", Status: 401, Title: "Forbidden", Detail: "This operation is forbidden."}
)

func newContentNotFoundError(contentName string) *api.Error {
	return &api.Error{ID: "not_found", Status: 404, Title: "Not Found", Detail: fmt.Sprintf("%s not found.", contentName)}
}

func newContentNotFoundErrorWithMessage(message string) *api.Error {
	return &api.Error{ID: "not_found", Status: 404, Title: "Not Found", Detail: message}
}

func newInternalServerError(err interface{}) *api.Error {
	return &api.Error{ID: "internal_server_error", Status: 500, Title: "Internal Server Error", Detail: fmt.Sprintf("Something went wrong: %+v", err)}
}

func newNotAcceptableError(accept string) *api.Error {
	return &api.Error{ID: "not_acceptable", Status: 406, Title: "Not Acceptable", Detail: fmt.Sprintf("Accept header must be set to '%s'.", accept)}
}

func newUnsupportedMediaTypeError(contentType string) *api.Error {
	return &api.Error{ID: "unsupported_media_type", Status: 415, Title: "Unsupported Media Type", Detail: fmt.Sprintf("Content-Type header must be set to: '%s'.", contentType)}
}

func newBadRequestParameterWithParam(param string, err error) *api.Error {
	return &api.Error{ID: "bad_request", Status: http.StatusBadRequest, Title: "Bad Request", Detail: fmt.Sprintf("Invalid %q parameter %v", param, err)}
}

func newBadRequestParameter(err error) *api.Error {
	return &api.Error{ID: "bad_request", Status: http.StatusBadRequest, Title: "Bad Request", Detail: fmt.Sprintf(err.Error())}
}

func newBadRequestError(err error) *api.Error {
	return &api.Error{ID: "bad_request", Status: http.StatusBadRequest, Title: "Bad Request", Detail: fmt.Sprint(err)}
}

func newBadRequestMessage(message string) *api.Error {
	return &api.Error{ID: "bad_request", Status: http.StatusBadRequest, Title: "Bad Request", Detail: message}
}

func newConflictRequest(message string) *api.Error {
	return &api.Error{ID: "conflict", Status: http.StatusConflict, Title: "Conflict", Detail: message}
}

func newForbiddenRequest(message string) *api.Error {
	return &api.Error{ID: "forbidden", Status: http.StatusForbidden, Title: "Forbidden", Detail: message}
}

func newUnauthorizedRequest(gc *gin.Context, realm string) *api.Error {
	gc.Header("WWW-Authenticate", realm)
	return &api.Error{ID: "unauthorized", Status: http.StatusUnauthorized, Title: "Unauthorized", Detail: realm}
}

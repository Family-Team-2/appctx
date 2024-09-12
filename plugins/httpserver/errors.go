package httpserver

import (
	"net/http"

	"github.com/go-chi/render"
)

type HTTPErrorer interface {
	Code() int
	Explanation() string
	Err() error
}

type HTTPGenericError struct {
	err  error
	expl string
	code int
}

var (
	httpErrorInternal         = HTTPGenericError{code: http.StatusInternalServerError, expl: "internal server error"}
	httpErrorMethodNotAllowed = HTTPGenericError{code: http.StatusMethodNotAllowed, expl: "method not allowed"}
	httpErrorNotFound         = HTTPGenericError{code: http.StatusNotFound, expl: "route not found"}
)

func (e *HTTPGenericError) Code() int {
	if e.code == 0 {
		return http.StatusInternalServerError
	}

	return e.code
}

func (e *HTTPGenericError) Explanation() string {
	if e.expl == "" {
		return "internal server error"
	}

	return e.expl
}

func (e *HTTPGenericError) Err() error {
	return e.err
}

func makeHTTPErrorer(e error) HTTPErrorer {
	// todo: handle friendly error (check with errors.As)
	return &HTTPGenericError{
		code: http.StatusInternalServerError,
		err:  e,
		expl: "internal server error",
	}
}

func (pl *PluginHTTPServer[T, U]) SendHTTPError(w http.ResponseWriter, r *http.Request, e HTTPErrorer) {
	if r.Context().Err() == nil {
		// send error to client only if context is not yet cancelled

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Code())
		render.DefaultResponder(w, r, render.M{
			"ok":    false,
			"error": e.Explanation(),
		})
	}
}

func (pl *PluginHTTPServer[T, U]) ThrowHTTPError(w http.ResponseWriter, r *http.Request, e HTTPErrorer) {
	pl.SendHTTPError(w, r, e)
	panic(e)
}

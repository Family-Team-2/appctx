package httpserver

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func (pl *PluginHTTPServer[T, U]) middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()

		defer func() {
			pl.app.Log().
				Str("method", r.Method).
				Str("proto", r.Proto).
				Str("path", r.URL.Path).
				Str("host", r.Host).
				Str("addr", r.RemoteAddr).
				Dur("time", time.Since(t1)).
				Int("status", ww.Status()).
				Int("size", ww.BytesWritten()).
				Msg("served request")
		}()

		next.ServeHTTP(ww, r)
	})
}

func (pl *PluginHTTPServer[T, U]) middlewareRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			errorer, ok := e.(HTTPErrorer)
			if ok {
				pl.app.Warn().Err(errorer.Err()).Str("explanation", errorer.Explanation()).Msg("error in http handlers")
				return
			}

			switch e := e.(type) {
			case error:
				if errors.Is(e, http.ErrAbortHandler) {
					panic(e)
				}

				pl.app.Warn().Err(e).Msg("error in http handlers")
			case string:
				pl.app.Warn().Str("error", e).Msg("error in http handlers")
			default:
				pl.app.Warn().Any("error", e).Msg("error in http handlers")
			}

			pl.SendHTTPError(w, r, &httpErrorInternal)
		}()

		next.ServeHTTP(w, r)
	})
}

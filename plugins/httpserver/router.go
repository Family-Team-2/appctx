package httpserver

import (
	"io/fs"
	"net/http"

	"github.com/Family-Team-2/appctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (pl *PluginHTTPServer[T, U]) buildRouter() {
	pl.rt = chi.NewRouter()
	pl.rt.Use(
		middleware.RealIP,
		middleware.CleanPath,
	)

	if pl.LogRequests {
		pl.rt.Use(pl.middlewareLogger)
	}

	pl.rt.Use(pl.middlewareRecoverer)

	pl.rt.NotFound(func(w http.ResponseWriter, r *http.Request) {
		pl.SendHTTPError(w, r, &httpErrorNotFound)
	})

	pl.rt.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		pl.SendHTTPError(w, r, &httpErrorMethodNotAllowed)
	})
}

type HTTPRouter[T any, U any] interface {
	Route(method, pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error)
	Subroute(pattern string, callback func(sr HTTPRouter[T, U]))

	Get(pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error)
	Post(pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error)
	Patch(pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error)
	Delete(pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error)
	ServeFS(pattern string, fs fs.FS)

	Use(middlewares ...func(app *appctx.AppCtx[T, U], next http.Handler) http.Handler)
	Group(callback func(sr HTTPRouter[T, U]))
}

type httpRouter[T any, U any] struct {
	r  chi.Router
	pl *PluginHTTPServer[T, U]
}

func (rt *httpRouter[T, U]) Route(method, pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error) {
	rt.r.MethodFunc(method, pattern, func(w http.ResponseWriter, r *http.Request) {
		err := handler(rt.pl.app, w, r)
		if err != nil {
			rt.pl.ThrowHTTPError(w, r, makeHTTPErrorer(err))
		}
	})
}

func (rt *httpRouter[T, U]) Get(pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error) {
	rt.Route("GET", pattern, handler)
}

func (rt *httpRouter[T, U]) Post(pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error) {
	rt.Route("POST", pattern, handler)
}

func (rt *httpRouter[T, U]) Patch(pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error) {
	rt.Route("PATCH", pattern, handler)
}

func (rt *httpRouter[T, U]) Delete(pattern string, handler func(app *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error) {
	rt.Route("DELETE", pattern, handler)
}

func (rt *httpRouter[T, U]) ServeFS(pattern string, fs fs.FS) {
	rt.Get(pattern, func(_ *appctx.AppCtx[T, U], w http.ResponseWriter, r *http.Request) error {
		http.FileServer(http.FS(fs)).ServeHTTP(w, r)
		return nil
	})
}

func (rt *httpRouter[T, U]) Use(middlewares ...func(app *appctx.AppCtx[T, U], next http.Handler) http.Handler) {
	for _, middleware := range middlewares {
		rt.r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := middleware(rt.pl.app, next)
				if h != nil {
					h.ServeHTTP(w, r)
				}
			})
		})
	}
}

func (rt *httpRouter[T, U]) Group(callback func(sr HTTPRouter[T, U])) {
	rt.r.Group(func(r chi.Router) {
		callback(&httpRouter[T, U]{
			r:  r,
			pl: rt.pl,
		})
	})
}

func (rt *httpRouter[T, U]) Subroute(pattern string, callback func(sr HTTPRouter[T, U])) {
	rt.r.Route(pattern, func(r chi.Router) {
		callback(&httpRouter[T, U]{
			r:  r,
			pl: rt.pl,
		})
	})
}

func (pl *PluginHTTPServer[T, U]) DefineHTTPRoutes(callback func(r HTTPRouter[T, U])) {
	callback(&httpRouter[T, U]{
		r:  pl.rt,
		pl: pl,
	})
}

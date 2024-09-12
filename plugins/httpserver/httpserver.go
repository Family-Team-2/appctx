package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Family-Team-2/appctx"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type PluginHTTPServer[T any, U any] struct {
	Host              string        `yaml:"host"`
	Port              uint16        `yaml:"port"`
	LogRequests       bool          `yaml:"log_requests"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`

	srv *http.Server
	rt  *chi.Mux
	app *appctx.AppCtx[T, U]
}

func (pl *PluginHTTPServer[T, U]) MarshalZerologObject(e *zerolog.Event) {
	e.Str("host", pl.Host).Uint16("port", pl.Port)
}

func (pl *PluginHTTPServer[T, U]) PluginName() string {
	return "httpserver"
}

func (pl *PluginHTTPServer[T, U]) PluginInstantiate(app *appctx.AppCtx[T, U]) error {
	pl.app = app

	pl.Host = "0.0.0.0"
	pl.Port = 80
	pl.ReadHeaderTimeout = 1 * time.Minute
	pl.ShutdownTimeout = 5 * time.Second
	return nil
}

func (pl *PluginHTTPServer[T, U]) PluginStart(app *appctx.AppCtx[T, U]) error {
	pl.buildRouter()

	pl.srv = &http.Server{
		Addr:    pl.Host + ":" + strconv.FormatUint(uint64(pl.Port), 10),
		Handler: pl.rt,
		BaseContext: func(_ net.Listener) context.Context {
			return app
		},
		ReadTimeout:       pl.ReadTimeout,
		ReadHeaderTimeout: pl.ReadHeaderTimeout,
	}

	return nil
}

func (pl *PluginHTTPServer[T, U]) StartServer() {
	pl.app.Debug().EmbedObject(pl).Msg("starting http server")

	go func() {
		err := pl.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pl.app.Error(err).Msg("failed to start http server")
			pl.app.Stop()
		}
	}()
}

func (pl *PluginHTTPServer[T, U]) StopServer() {
	pl.app.Debug().EmbedObject(pl).Msg("stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), pl.ShutdownTimeout)
	defer cancel()

	err := pl.srv.Shutdown(ctx)
	if err != nil {
		pl.app.Warn().Err(err).Msg("failed to gracefully shutdown http server")
		err := pl.srv.Close()
		if err != nil {
			pl.app.Warn().Err(err).Msg("failed to close http server")
		}
	}
}

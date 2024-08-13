package appctx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

type appCfg[T any, U any] struct {
	Debug      bool   `yaml:"debug"`
	ConfigFile string `yaml:"-"`

	Custom  T `yaml:",inline"`
	Plugins U `yaml:",inline"`
}

type AppCtx[T any, U any] struct {
	context.Context

	cfg    appCfg[T, U]
	logger zerolog.Logger

	title             string
	version           string
	hasLogger         bool
	cancel            func()
	flags             []appFlag
	registeredPlugins []AppPlugin[T, U]
}

func NewApp[T any, U any](title, version string) *AppCtx[T, U] {
	return &AppCtx[T, U]{
		title:   title,
		version: version,
	}
}

func (app *AppCtx[T, U]) Config() *T {
	return &app.cfg.Custom
}

func (app *AppCtx[T, U]) Plugins() *U {
	return &app.cfg.Plugins
}

func (app *AppCtx[T, U]) C() *T {
	return app.Config()
}

func (app *AppCtx[T, U]) P() *U {
	return app.Plugins()
}

func (app *AppCtx[T, U]) Run(callback func(ctx *AppCtx[T, U]) error) {
	app.Context, app.cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer app.cancel()

	defer app.stopPlugins()

	err := app.run(callback)
	if err != nil {
		if app.hasLogger {
			app.logger.Err(err).Msg("shutting down")
		} else {
			fmt.Println("ERROR: " + err.Error())
		}
	} else {
		app.logger.Info().Msg("shutting down")
	}
}

func (app *AppCtx[_, _]) Stop() {
	app.logger.Debug().Msg("app stop requested")

	if app.cancel != nil {
		app.cancel()
	}
}

func (app *AppCtx[T, U]) run(callback func(ctx *AppCtx[T, U]) error) error {
	setDefault(&app.title, "App")
	setDefault(&app.version, "0.0.1")

	app.Flag2("d", "debug", &app.cfg.Debug, false, "enable debug output")
	app.Flag2("c", "config-file", &app.cfg.ConfigFile, "config.yml", "path to config file")

	err := app.initFlags()
	if err != nil {
		return fmt.Errorf("initializing flags: %w", err)
	}

	err = app.loadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	app.makeLogger()

	err = app.startPlugins()
	if err != nil {
		return fmt.Errorf("starting plugins: %w", err)
	}

	app.Log().Str("title", app.title).Str("version", app.version).Msg("app: running")
	return callback(app)
}

func (app *AppCtx[T, U]) clone() *AppCtx[T, U] {
	newApp := *app
	return &newApp
}

func setDefault[T comparable](v *T, def T) {
	var zero T
	if *v == zero {
		*v = def
	}
}

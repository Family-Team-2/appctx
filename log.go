package appctx

import (
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func (app *AppCtx[T, U]) Logger() *zerolog.Logger {
	return &app.logger
}

func (app *AppCtx[T, U]) Log() *zerolog.Event {
	return app.logger.Info() //nolint:zerologlint
}

func (app *AppCtx[T, U]) Warn() *zerolog.Event {
	return app.logger.Warn() //nolint:zerologlint
}

func (app *AppCtx[T, U]) Debug() *zerolog.Event {
	return app.logger.Debug() //nolint:zerologlint
}

func (app *AppCtx[T, U]) Error(errs ...error) *zerolog.Event {
	if len(errs) == 0 {
		return app.logger.Error() //nolint:zerologlint
	}

	return app.logger.Error().Err(errors.Join(errs...)) //nolint:zerologlint
}

func (app *AppCtx[_, _]) makeLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	app.logger = zerolog.New(os.Stdout)
	if app.cfg.Debug {
		app.logger = app.logger.Level(zerolog.DebugLevel).Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "02.01.2006 15:04:05.000000",
		})
	} else {
		app.logger = app.logger.Level(zerolog.InfoLevel)
	}

	app.logger = app.logger.With().Timestamp().Logger()
	app.logger.Debug().Msg("logger: initialized")
	app.hasLogger = true
}

package appctx

import (
	"errors"
	"fmt"
	"slices"
)

type AppPlugin[T any, U any] interface {
	PluginName() string
}

func (app *AppCtx[T, U]) RegisterPlugin(plugin AppPlugin[T, U]) {
	app.registeredPlugins = append(app.registeredPlugins, plugin)
}

func (app *AppCtx[T, U]) WithPlugin(plugin AppPlugin[T, U]) *AppCtx[T, U] {
	app.RegisterPlugin(plugin)
	return app
}

type appPluginStarter[T any, U any] interface {
	AppPlugin[T, U]
	PluginStart(ac *AppCtx[T, U]) error
}

type appPluginStopper[T any, U any] interface {
	AppPlugin[T, U]
	PluginStop(ac *AppCtx[T, U])
}

func (app *AppCtx[T, U]) startPlugins() error {
	errs := []error{}

	for _, rawPlugin := range app.registeredPlugins {
		plugin, ok := rawPlugin.(appPluginStarter[T, U])
		if !ok {
			continue
		}

		name := plugin.PluginName()
		app.Debug().Str("name", name).Msg("starting plugin")

		err := plugin.PluginStart(app)
		if err != nil {
			errs = append(errs, fmt.Errorf("starting plugin \"%v\": %w", name, err))
		}
	}

	return errors.Join(errs...)
}

func (app *AppCtx[T, U]) stopPlugins() {
	plugins := app.registeredPlugins
	slices.Reverse(plugins)

	for _, rawPlugin := range plugins {
		plugin, ok := rawPlugin.(appPluginStopper[T, U])
		if !ok {
			continue
		}

		app.Debug().Str("name", plugin.PluginName()).Msg("stopping plugin")

		plugin.PluginStop(app)
	}
}

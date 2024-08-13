package main

import (
	"github.com/Family-Team-2/appctx"
	"github.com/Family-Team-2/appctx/plugins/example"
)

func main() {
	type appConfig struct {
		Message string
	}

	type appPlugins struct {
		example.PluginExample[appConfig, appPlugins] `yaml:",inline"`
	}

	type App = appctx.AppCtx[appConfig, appPlugins]

	app := appctx.NewApp[appConfig, appPlugins]("Example App", "1.0.0")
	app.RegisterPlugin(&app.Plugins().PluginExample)

	app.Run(func(app *App) error {
		app.Log().Str("config_message", app.Config().Message).Msg("running app callback")
		app.Log().Str("string", app.Plugins().PluginTestFunction()).Msg("calling plugin function")
		return nil
	})
}

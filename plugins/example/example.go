// Package example represents a barebone plugin that can be used as an example.

// This plugin implements Start/Stop callbacks.
// It also exports a single YAML string which gets filled in
// during the configuration parsing stage.

package example

import (
	"github.com/Family-Team-2/appctx"
)

// PluginExample is an example plugin for AppCtx.
type PluginExample[T any, U any] struct {
	PluginConfigItem string `yaml:"plugin_config_item"`
}

func (pl *PluginExample[T, U]) PluginName() string {
	return "example"
}

func (pl *PluginExample[T, U]) PluginStart(app *appctx.AppCtx[T, U]) error {
	app.Log().Msg("example: running PluginStart()")
	return nil
}

func (pl *PluginExample[T, U]) PluginStop(app *appctx.AppCtx[T, U]) {
	app.Log().Msg("example: running PluginStop()")
}

// PluginFunction gets exported to AppCtx and can be called from App.Run() method.
func (pl *PluginExample[T, U]) PluginTestFunction() string {
	return "Example plugin: value of PluginConfigItem is \"" + pl.PluginConfigItem + "\""
}

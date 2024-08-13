package appctx

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetCommandlineFlags() {
	os.Args = []string{os.Args[0]}
}

func TestSimpleApp(t *testing.T) {
	resetCommandlineFlags()

	f := false

	app := NewApp[struct{}, struct{}]("Test App", "1.0.0")
	app.DisableConfig()
	app.Run(func(_ *AppCtx[struct{}, struct{}]) error {
		f = true
		return nil
	})

	assert.EqualValues(t, "Test App", app.title)
	assert.EqualValues(t, "1.0.0", app.version)
	require.True(t, f)
}

func TestAppErrorHandling(t *testing.T) {
	resetCommandlineFlags()

	app := NewApp[struct{}, struct{}]("Test App", "1.0.0")
	app.DisableConfig()
	app.Run(func(_ *AppCtx[struct{}, struct{}]) error {
		return assert.AnError
	})

	require.True(t, app.hasError)
}

func TestAppStop(t *testing.T) {
	resetCommandlineFlags()

	app := NewApp[struct{}, struct{}]("Test App", "1.0.0")
	app.DisableConfig()
	app.Run(func(app *AppCtx[struct{}, struct{}]) error {
		app.Stop()

		<-app.Context.Done()

		assert.ErrorIs(t, context.Canceled, app.Context.Err())
		return app.Context.Err()
	})

	require.True(t, app.hasError)
}

func TestAppConfig(t *testing.T) {
	resetCommandlineFlags()

	type appConfig struct {
		TestStr string
	}

	app := NewApp[appConfig, struct{}]("Test App", "1.0.0")
	app.C().TestStr = "1"
	app.DisableConfig()
	app.Run(func(app *AppCtx[appConfig, struct{}]) error {
		assert.EqualValues(t, app.C().TestStr, "1")
		return nil
	})

	require.False(t, app.hasError)
}

type appTestPlugin[T any, U any] struct {
	TestStr string
	started bool
	stopped bool
}

func (pl *appTestPlugin[T, U]) PluginName() string {
	return "testing"
}

func (pl *appTestPlugin[T, U]) PluginStart(_ *AppCtx[T, U]) error {
	pl.started = true
	return nil
}

func (pl *appTestPlugin[T, U]) PluginStop(_ *AppCtx[T, U]) {
	pl.stopped = true
}

func TestAppPlugin(t *testing.T) {
	resetCommandlineFlags()

	type appPlugins struct {
		appTestPlugin[struct{}, appPlugins] `yaml:",inline"`
	}

	app := NewApp[struct{}, appPlugins]("Test App", "1.0.0")
	app.P().TestStr = "1"
	app.RegisterPlugin(&app.P().appTestPlugin)
	app.DisableConfig()
	app.Run(func(app *AppCtx[struct{}, appPlugins]) error {
		assert.EqualValues(t, app.P().TestStr, "1")

		assert.True(t, app.P().started)
		assert.False(t, app.P().stopped)

		return nil
	})

	assert.True(t, app.P().started)
	assert.True(t, app.P().stopped)
	require.False(t, app.hasError)
}

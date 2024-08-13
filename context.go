package appctx

import (
	"context"
	"time"
)

func (app *AppCtx[T, U]) WithValue(key, val any) *AppCtx[T, U] {
	return app.cloneWithContext(context.WithValue(app.Context, key, val))
}

func (app *AppCtx[T, U]) WithTimeout(d time.Duration) (newApp *AppCtx[T, U], done func()) {
	ctx, cancel := context.WithTimeout(app.Context, d)
	return app.cloneWithContext(ctx), cancel
}

func (app *AppCtx[T, U]) WithCancel() (newApp *AppCtx[T, U], done func()) {
	ctx, cancel := context.WithCancel(app.Context)
	return app.cloneWithContext(ctx), cancel
}

func (app *AppCtx[T, U]) WithContext(ctx context.Context) *AppCtx[T, U] {
	return app.cloneWithContext(ctx)
}

func (app *AppCtx[T, U]) cloneWithContext(ctx context.Context) *AppCtx[T, U] {
	newApp := app.clone()
	newApp.Context = ctx
	return newApp
}

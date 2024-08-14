package gorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Family-Team-2/appctx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PluginGORM[T any, U any] struct {
	DatabaseURL           string        `yaml:"database_url"`
	TraceSQL              bool          `yaml:"trace_sql"`
	MaxConnectionLifetime time.Duration `yaml:"max_connection_lifetime"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`

	db *gorm.DB
}

func (pl *PluginGORM[T, U]) PluginName() string {
	return "gorm"
}

func (pl *PluginGORM[T, U]) PluginInstantiate(_ *appctx.AppCtx[T, U]) error {
	pl.MaxConnectionLifetime = 5 * time.Minute
	pl.MaxOpenConnections = 10
	return nil
}

func (pl *PluginGORM[T, U]) PluginStart(app *appctx.AppCtx[T, U]) error {
	if pl.DatabaseURL == "" {
		return errors.New("empty database URL")
	}

	db, err := gorm.Open(postgres.Open(pl.DatabaseURL), &gorm.Config{
		Logger: NewLogger(app.Logger(), pl.TraceSQL),
	})
	if err != nil {
		return fmt.Errorf("initializing db: %w", err)
	}

	pl.db = db.WithContext(app)

	sqlDB, err := pl.sqlDB()
	if err != nil {
		return fmt.Errorf("getting SQL DB: %w", err)
	}

	sqlDB.SetConnMaxLifetime(pl.MaxConnectionLifetime)
	sqlDB.SetMaxOpenConns(pl.MaxOpenConnections)

	err = pl.testDBFeatures(app)
	if err != nil {
		return fmt.Errorf("testing db features: %w", err)
	}

	return nil
}

func (pl *PluginGORM[T, U]) PluginStop(app *appctx.AppCtx[T, U]) {
	sqlDB, err := pl.sqlDB()
	if err != nil {
		app.Error(err).Msg("failed to get sql DB instance")
		return
	}

	err = sqlDB.Close()
	if err != nil {
		app.Error(err).Msg("failed to close sql DB")
	}

	pl.db = nil
}

func (pl *PluginGORM[T, U]) DB() *gorm.DB {
	return pl.db
}

func (pl *PluginGORM[T, U]) DBC(ctx context.Context) *gorm.DB {
	return pl.db.WithContext(ctx)
}

func (pl *PluginGORM[T, U]) sqlDB() (*sql.DB, error) {
	sqlDB, err := pl.db.DB()
	if err != nil {
		return nil, err
	}

	return sqlDB, nil
}

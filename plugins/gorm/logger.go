package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type DBLogger struct {
	logger *zerolog.Logger
	trace  bool
}

func NewLogger(logger *zerolog.Logger, trace bool) *DBLogger {
	return &DBLogger{
		logger: logger,
		trace:  trace,
	}
}

func (l DBLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l DBLogger) Error(_ context.Context, msg string, opts ...any) {
	l.logger.Error().Msg(fmt.Sprintf(msg, opts...))
}

func (l DBLogger) Warn(_ context.Context, msg string, opts ...any) {
	l.logger.Warn().Msg(fmt.Sprintf(msg, opts...))
}

func (l DBLogger) Info(_ context.Context, msg string, opts ...any) {
	l.logger.Info().Msg(fmt.Sprintf(msg, opts...))
}

func (l DBLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, _ int64), _ error) {
	if !l.trace {
		return
	}

	event := l.logger.Info().Dur("elapsed", time.Since(begin))

	sql, rows := fc()
	if sql != "" {
		event = event.Str("sql", sql)
	}
	if rows > -1 {
		event = event.Int64("rows", rows)
	}

	event.Send()
}

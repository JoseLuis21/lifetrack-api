package logging

import (
	"context"
	"log"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewZap(lc fx.Lifecycle) *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("can't initialize zap logger %+v", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return err
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing zap logger")
			return logger.Sync()
		},
	})

	return logger
}

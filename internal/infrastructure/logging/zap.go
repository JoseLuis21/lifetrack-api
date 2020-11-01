package logging

import (
	"context"
	"log"

	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewZap creates a preconfigured zap.Logger
func NewZap(lc fx.Lifecycle, cfg configuration.Configuration) *zap.Logger {
	logger, err := getZapLogger(cfg)
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

func getZapLogger(cfg configuration.Configuration) (*zap.Logger, error) {
	if cfg.IsDevEnv() || cfg.IsTestEnv() {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}

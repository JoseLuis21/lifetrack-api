package main

import (
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/lifetrack-api/internal/application/category"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/logging"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence/cassandracategory"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence/dbpool"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/categoryhandler"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/serverless"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			configuration.NewConfiguration,
			logging.NewZap,
			dbpool.NewCassandra,
			func(logger *zap.Logger, s *gocql.Session) repository.Category {
				return persistence.NewCategory(cassandracategory.NewRepository(s), logger)
			},
			func(cfg configuration.Configuration, logger *zap.Logger) event.Bus {
				return eventbus.New(eventbus.NewInMemory(cfg), cfg, logger)
			},
			transport.NewMux,
			category.NewAddCommandHandler,
			func(cmd *category.AddCommandHandler, r *mux.Router) serverless.Handler {
				return categoryhandler.NewAdd(cmd, r)
			},
		),
		fx.Invoke(serverless.StartLambda),
	)
	app.Run()
	select {
	case <-app.Done():
	}
}

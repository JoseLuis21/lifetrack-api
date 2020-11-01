package main

import (
	"github.com/neutrinocorp/lifetrack-api/internal/application/category"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/logging"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence/inmemcategory"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/categoryhandler"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			configuration.NewConfiguration,
			logging.NewZap,
			func(logger *zap.Logger) repository.Category {
				return persistence.NewCategory(inmemcategory.NewInMemory(), logger)
			},
			func(cfg configuration.Configuration, logger *zap.Logger) event.Bus {
				return eventbus.New(eventbus.NewInMemory(cfg), cfg, logger)
			},
			transport.NewMux,
			category.NewGetQuery,
			category.NewAddCommandHandler,
			category.NewListQuery,
			category.NewUpdateCommandHandler,
			category.NewRemoveCommandHandler,
		),
		fx.Invoke(
			categoryhandler.NewGet,
			categoryhandler.NewList,
			categoryhandler.NewAdd,
			categoryhandler.NewEdit,
			categoryhandler.NewRemove,
		),
	)
	app.Run()

	select {
	case <-app.Done():
	}
}

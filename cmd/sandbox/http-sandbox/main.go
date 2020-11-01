package main

import (
	"context"
	"log"
	"net/http"

	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/logging"
	"go.uber.org/zap"

	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence"

	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence/inmemcategory"

	"github.com/neutrinocorp/lifetrack-api/internal/application/category"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/categoryhandler"
	"go.uber.org/fx"

	"github.com/gorilla/mux"
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
			NewMux,
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
		log.Print("stop")
	}
}

func NewMux(lc fx.Lifecycle) *mux.Router {
	r := mux.NewRouter()
	r = r.PathPrefix("/live").Subrouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Print("starting http server")
				_ = server.ListenAndServe()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Print("stopping http server")
			return server.Shutdown(ctx)
		},
	})

	return r
}

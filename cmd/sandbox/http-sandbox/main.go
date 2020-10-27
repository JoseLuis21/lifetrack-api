package main

import (
	"context"
	"log"
	"net/http"

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
	// logger, _ := zap.NewProduction()
	app := fx.New(
		fx.Provide(
			func() repository.Category {
				return inmemcategory.NewInMemory()
			},
			configuration.NewConfiguration,
			func(cfg configuration.Configuration) event.Bus {
				return eventbus.NewInMemory(cfg)
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
	// Add middlewares

	// Known code-smell, ignore
	/*
		getCategory, cleanCGet, err := dep.InjectGetCategoryQuery()
		if err != nil {
			panic(err)
		}
		defer cleanCGet()

		_ = categoryhandler.NewGet(getCategory, r)

		listCategory, cleanLCat, err := dep.InjectListCategoriesQuery()
		if err != nil {
			panic(err)
		}
		defer cleanLCat()

		_ = categoryhandler.NewList(listCategory, r)

		addCategory, cleanACat, err := dep.InjectAddCategoryHandler()
		if err != nil {
			panic(err)
		}
		defer cleanACat()

		_ = categoryhandler.NewAdd(addCategory, r)

		editCategory, cleanECat, err := dep.InjectEditCategory()
		if err != nil {
			panic(err)
		}
		defer cleanECat()

		_ = categoryhandler.NewEdit(editCategory, r)

		removeCategory, cleanRCat, err := dep.InjectRemoveCategory()
		if err != nil {
			panic(err)
		}
		defer cleanRCat()

		_ = categoryhandler.NewRemove(removeCategory, r)

		log.Print("starting http sandbox")

		panic(http.ListenAndServe(":8080", r))*/
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

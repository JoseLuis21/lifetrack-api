//+build wireinject

package dep

import (
	"github.com/google/wire"
	"github.com/neutrinocorp/life-track-api/internal/application/category"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/inmemcategory"
)

var infraSet = wire.NewSet(
	configuration.NewConfiguration,
	// logging.NewZapProd,
	wire.Bind(new(repository.Category), new(*inmemcategory.InMemory)),
	inmemcategory.NewInMemory,
	provideEventBus,
	eventbus.NewInMemory,
)

func provideEventBus(cfg configuration.Configuration) event.Bus {
	return eventbus.NewInMemory(cfg)
}

func InjectAddCategoryHandler() (*category.AddCommandHandler, func(), error) {
	wire.Build(infraSet, category.NewAddCommandHandler)
	return &category.AddCommandHandler{}, nil, nil
}

func InjectGetCategoryQuery() (*category.GetQuery, func(), error) {
	wire.Build(infraSet, category.NewGetQuery)
	return &category.GetQuery{}, nil, nil
}

func InjectListCategoriesQuery() (*category.ListQuery, func(), error) {
	wire.Build(infraSet, category.NewListQuery)
	return &category.ListQuery{}, nil, nil
}

func InjectEditCategory() (*category.UpdateCommandHandler, func(), error) {
	wire.Build(infraSet, category.NewUpdateCommandHandler)
	return &category.UpdateCommandHandler{}, nil, nil
}

func InjectRemoveCategory() (*category.RemoveCommandHandler, func(), error) {
	wire.Build(infraSet, category.NewRemoveCommandHandler)
	return &category.RemoveCommandHandler{}, nil, nil
}

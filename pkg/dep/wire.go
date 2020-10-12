//+build wireinject

package dep

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/wire"
	categoryapp "github.com/neutrinocorp/life-track-api/internal/application/category"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/awsutil"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/logging"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/category"
	"go.uber.org/zap"
)

var infraSet = wire.NewSet(
	infrastructure.NewConfiguration,
	awsutil.NewSession,
	logging.NewZapProd,
	provideCategoryRepository,
	provideEventBus,
	eventbus.NewAWS,
)

func provideCategoryRepository(s *session.Session, cfg infrastructure.Configuration, logger *zap.Logger) repository.Category {
	return category.NewCategory(category.NewDynamoRepository(s, cfg), logger)
}

func provideEventBus(s *session.Session, cfg infrastructure.Configuration, logger *zap.Logger) event.Bus {
	return eventbus.NewEventBus(eventbus.NewAWS(s, cfg), logger)
}

func InjectAddCategoryHandler() (*categoryapp.AddHandler, func(), error) {
	wire.Build(infraSet, categoryapp.NewAddHandler)
	return &categoryapp.AddHandler{}, nil, nil
}

func InjectGetCategoryQuery() (*categoryapp.Get, func(), error) {
	wire.Build(infraSet, categoryapp.NewGet)
	return &categoryapp.Get{}, nil, nil
}

func InjectListCategoriesQuery() (*categoryapp.List, func(), error) {
	wire.Build(infraSet, categoryapp.NewList)
	return &categoryapp.List{}, nil, nil
}

func InjectChangeCategoryState() (*categoryapp.ChangeStateHandler, func(), error) {
	wire.Build(infraSet, categoryapp.NewChangeStateHandler)
	return &categoryapp.ChangeStateHandler{}, nil, nil
}

func InjectEditCategory() (*categoryapp.EditHandler, func(), error) {
	wire.Build(infraSet, categoryapp.NewEditHandler)
	return &categoryapp.EditHandler{}, nil, nil
}

func InjectRemoveCategory() (*categoryapp.RemoveHandler, func(), error) {
	wire.Build(infraSet, categoryapp.NewRemoveHandler)
	return &categoryapp.RemoveHandler{}, nil, nil
}

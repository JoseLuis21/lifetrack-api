//+build wireinject

package dep

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/wire"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
	"github.com/neutrinocorp/life-track-api/internal/application/query"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/logging"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
	"go.uber.org/zap"
)

var infraSet = wire.NewSet(
	infrastructure.NewConfiguration,
	infrastructure.NewSession,
	logging.NewZapProd,
	provideCategoryRepository,
	wire.Bind(new(event.Bus), new(*eventbus.AWS)),
	eventbus.NewAWS,
)

func provideCategoryRepository(s *session.Session, cfg infrastructure.Configuration, logger *zap.Logger) repository.Category {
	return persistence.NewCategory(persistence.NewCategoryDynamoRepository(s, cfg), logger)
}

func InjectAddCategoryHandler() (*command.AddCategoryHandler, func(), error) {
	wire.Build(infraSet, command.NewAddCategoryHandler)

	return &command.AddCategoryHandler{}, nil, nil
}

func InjectGetCategoryQuery() (*query.GetCategory, func(), error) {
	wire.Build(infraSet, query.NewGetCategory)

	return &query.GetCategory{}, nil, nil
}

func InjectListCategoriesQuery() (*query.ListCategories, func(), error) {
	wire.Build(infraSet, query.NewListCategories)

	return &query.ListCategories{}, nil, nil
}

func InjectChangeCategoryState() (*command.ChangeCategoryStateHandler, func(), error) {
	wire.Build(infraSet, command.NewChangeCategoryStateHandler)

	return &command.ChangeCategoryStateHandler{}, nil, nil
}

func InjectEditCategory() (*command.EditCategoryHandler, func(), error) {
	wire.Build(infraSet, command.NewEditCategoryHandler)

	return &command.EditCategoryHandler{}, nil, nil
}

func InjectRemoveCategory() (*command.RemoveCategoryHandler, func(), error) {
	wire.Build(infraSet, command.NewRemoveCategoryHandler)

	return &command.RemoveCategoryHandler{}, nil, nil
}

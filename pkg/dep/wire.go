//+build wireinject

package dep

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/wire"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
	"github.com/neutrinocorp/life-track-api/internal/application/query"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/logging"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
	"go.uber.org/zap"
)

var dynamoSet = wire.NewSet(
	infrastructure.NewConfiguration,
	infrastructure.NewSession,
	logging.NewZapProd,
	provideCategoryRepository,
)

func provideCategoryRepository(s *session.Session, cfg infrastructure.Configuration, logger *zap.Logger) repository.Category {
	return persistence.NewCategory(persistence.NewCategoryDynamoRepository(s, cfg), logger)
}

func InjectAddCategoryHandler() (*command.AddCategoryHandler, func(), error) {
	wire.Build(dynamoSet, command.NewAddCategoryHandler)

	return &command.AddCategoryHandler{}, nil, nil
}

func InjectGetCategoryQuery() (*query.GetCategory, func(), error) {
	wire.Build(dynamoSet, query.NewGetCategory)

	return &query.GetCategory{}, nil, nil
}

func InjectListCategoriesQuery() (*query.ListCategories, func(), error) {
	wire.Build(dynamoSet, query.NewListCategories)

	return &query.ListCategories{}, nil, nil
}

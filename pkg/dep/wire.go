//+build wireinject

package dep

import (
	"github.com/google/wire"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
	"github.com/neutrinocorp/life-track-api/internal/application/query"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
)

var dynamoSet = wire.NewSet(
	infrastructure.NewConfiguration,
	infrastructure.NewSession,
	wire.Bind(new(repository.Category), new(*persistence.CategoryDynamoRepository)),
	persistence.NewCategoryDynamoRepository,
)

func InjectAddCategoryHandler() (*command.AddCategoryHandler, error) {
	wire.Build(dynamoSet, command.NewAddCategoryHandler)

	return &command.AddCategoryHandler{}, nil
}

func InjectGetCategoryQuery() (*query.GetCategory, error) {
	wire.Build(dynamoSet, query.NewGetCategory)

	return &query.GetCategory{}, nil
}

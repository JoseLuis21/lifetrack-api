package main

import (
	"context"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/application/query"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
	"log"
)

func main() {
	cfg, err := infrastructure.NewConfiguration()
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	q := query.NewGetCategory(persistence.NewCategoryDynamoRepository(cfg))

	category, err := q.Query(context.Background(), "34c71dff-fcc4-46ce-bebb-8ab71bf45edb")
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	log.Print(category)
}

package main

import (
	"context"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/aggregate"
	"log"

	"github.com/gocql/gocql"

	"go.uber.org/fx"

	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/logging"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence/dbpool"
)

func main() {
	app := fx.New(
		fx.Provide(
			configuration.NewConfiguration,
			logging.NewZap,
			dbpool.NewCassandra,
		),
		fx.Invoke(run),
	)

	err := app.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = app.Stop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func run(s *gocql.Session) {
	log.Print("on workflow")

	c, err := aggregate.NewCategory("aruiz", "Gaming")
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}
	_ = c.ModifyDescription("Gaming")

	if err = s.Query(`INSERT INTO category (id, user_id, name, description,
		    target_time, picture, create_time, update_time, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID(), c.User(), c.Name(), c.Description(), int64(c.TargetTime().Minutes()), c.Picture(), c.CreateTime(), c.UpdateTime(), c.State()).Exec(); err != nil {
		log.Fatal(err)
	}
}

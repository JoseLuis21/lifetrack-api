package main

import (
	"log"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/google/uuid"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/dbpool"
)

func main() {
	cfg := configuration.NewConfiguration()
	s, clean, err := dbpool.NewCassandra(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer clean()

	c, err := aggregate.NewCategory(uuid.New().String(), "Physics")
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}
	_ = c.ModifyDescription("This is a sample")
	// c := cAg.MarshalPrimitive()

	if err = s.Query(`INSERT INTO category (id, user_id, name, description,
    target_time, picture, create_time, update_time, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID(), c.User(), c.Name(), c.Description(), int64(c.TargetTime().Minutes()), c.Picture(), c.CreateTime(), c.UpdateTime(), c.State()).Exec(); err != nil {
		log.Fatal(err)
	}
}

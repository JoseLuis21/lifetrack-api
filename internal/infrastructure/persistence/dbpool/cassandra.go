package dbpool

import (
	"context"
	"log"

	"github.com/gocql/gocql"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"go.uber.org/fx"
)

// NewCassandra creates an Apache Cassandra session
func NewCassandra(lc fx.Lifecycle, cfg configuration.Configuration) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Cassandra.Cluster...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Cassandra.Username,
		Password: cfg.Cassandra.Password,
	}
	cluster.Keyspace = cfg.Cassandra.Keyspace
	cluster.Consistency = gocql.LocalOne
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Print("closing cassandra session")
			session.Close()
			return nil
		},
	})

	return session, nil
}

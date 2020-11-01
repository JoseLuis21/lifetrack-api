package dbpool

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewCassandra creates an Apache Cassandra session
func NewCassandra(lc fx.Lifecycle, logger *zap.Logger, cfg configuration.Configuration) (*gocql.Session, error) {
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
			logger.Info("closing cassandra session",
				zap.String("ip_address", cfg.Cassandra.Cluster[0]),
				zap.String("keyspace", cfg.Cassandra.Keyspace),
			)
			session.Close()
			return nil
		},
	})

	return session, nil
}

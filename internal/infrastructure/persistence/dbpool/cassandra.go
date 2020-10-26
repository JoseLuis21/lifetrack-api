package dbpool

import (
	"github.com/gocql/gocql"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
)

// NewCassandra creates an Apache Cassandra session
func NewCassandra(cfg configuration.Configuration) (*gocql.Session, func(), error) {
	cluster := gocql.NewCluster(cfg.Cassandra.Cluster...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Cassandra.Username,
		Password: cfg.Cassandra.Password,
	}
	cluster.Keyspace = cfg.Cassandra.Keyspace
	cluster.Consistency = gocql.LocalQuorum
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, nil, err
	}

	return session, session.Close, nil
}

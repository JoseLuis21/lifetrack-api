package dbpool

import (
	"github.com/gocql/gocql"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/configuration"
)

// NewCassandra creates an Apache Cassandra session
func NewCassandra(cfg configuration.Configuration) (*gocql.Session, func(), error) {
	cluster := gocql.NewCluster(cfg.Cassandra.Cluster...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}
	cluster.Keyspace = cfg.Cassandra.Keyspace
	cluster.Consistency = gocql.Any
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		session.Close()
	}

	return session, cleanup, nil
}

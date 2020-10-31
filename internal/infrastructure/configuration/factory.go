package configuration

import (
	"github.com/spf13/viper"
)

// NewConfiguration reads and returns a kernel configuration
func NewConfiguration() Configuration {
	cfg := &Configuration{
		Cassandra:   &cassandra{},
		DynamoTable: &dynamoTable{},
	}
	if err := cfg.LoadEnv(); err != nil {
		return getDefault()
	}
	cfg.Read()

	return *cfg
}

func getDefault() Configuration {
	stage := viper.GetString("stage")
	return Configuration{
		Version: viper.GetString("version"),
		Stage:   stage,
		DynamoTable: &dynamoTable{
			Name:   setResourceStage(viper.GetString("dynamodb.table"), stage),
			Region: viper.GetString("dynamodb.region"),
		},
		Cassandra: &cassandra{
			Keyspace: setResourceStage(viper.GetString("cassandra.keyspace"), stage),
			Cluster:  viper.GetStringSlice("cassandra.cluster"),
			Username: viper.GetString("cassandra.username"),
			Password: viper.GetString("cassandra.password"),
		},
	}
}

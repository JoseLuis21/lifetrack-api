package configuration

import (
	"log"

	"github.com/spf13/viper"
)

// NewConfiguration reads and returns a kernel configuration
func NewConfiguration() Configuration {
	cfg := &Configuration{
		HTTP:        &httpServer{},
		Cassandra:   &cassandra{},
		DynamoTable: &dynamoTable{},
	}
	if err := cfg.LoadEnv(); err != nil {
		log.Printf("failed to load configuration: %+v", err)
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
		HTTP: &httpServer{
			Address:  viper.GetString("http.address"),
			Port:     viper.GetInt("http.port"),
			Endpoint: viper.GetString("http.endpoint"),
		},
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

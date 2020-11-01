package configuration

import "github.com/spf13/viper"

// cassandra Apache Cassandra config
type cassandra struct {
	Keyspace string   `json:"keyspace"`
	Cluster  []string `json:"cluster"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

func init() {
	viper.SetDefault("cassandra.keyspace", "lifetrack_"+DevelopmentStage)
	viper.SetDefault("cassandra.cluster", []string{"127.0.0.1"})

	// TODO: Migrate secrets to AWS Secrets Manager or Hashicorp Vault
	viper.SetDefault("cassandra.username", "cassandra")
	viper.SetDefault("cassandra.password", "cassandra")
}

func (c *cassandra) Load(stage string) {
	c.Keyspace = setResourceStage(viper.GetString("cassandra.keyspace"), stage)
	c.Cluster = viper.GetStringSlice("cassandra.cluster")
	c.Username = viper.GetString("cassandra.username")
	c.Password = viper.GetString("cassandra.password")
}

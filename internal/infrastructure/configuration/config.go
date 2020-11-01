package configuration

import (
	"os"

	"github.com/spf13/viper"
)

// Configuration kernel configuration
type Configuration struct {
	Version     string       `json:"version"`
	Stage       string       `json:"stage"`
	DynamoTable *dynamoTable `json:"dynamo_table"`
	Cassandra   *cassandra   `json:"cassandra"`
}

func init() {
	viper.SetDefault("version", "0.1.0-alpha")
	viper.SetDefault("stage", "prod")
}

func (c Configuration) LoadEnv() error {
	//	rule
	//	if stage is development (dev) then use config file (lifetrack.yaml)
	//	else use environment variables
	if c.isDevEnv() {
		return c.loadFromFile()
	}

	viper.SetEnvPrefix("lt")
	viper.AutomaticEnv()
	return nil
}

func (c Configuration) isDevEnv() bool {
	stage := os.Getenv("LT_STAGE")
	return stage == "dev"
}

func (c Configuration) loadFromFile() error {
	viper.SetConfigName("lifetrack")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.lifetrack/")
	viper.AddConfigPath("/etc/lifetrack/")
	viper.AddConfigPath("./conf/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return viper.SafeWriteConfig()
		}

		return err
	}

	viper.WatchConfig()
	return nil
}

func (c *Configuration) Read() {
	c.Version = viper.GetString("version")
	c.Stage = viper.GetString("stage")
	c.DynamoTable.Load(c.Stage)
	c.Cassandra.Load(c.Stage)
}
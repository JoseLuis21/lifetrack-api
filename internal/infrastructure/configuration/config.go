package configuration

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Configuration kernel configuration
type Configuration struct {
	Service     string       `json:"service"`
	Version     string       `json:"version"`
	Stage       string       `json:"stage"`
	HTTP        *httpServer  `json:"http"`
	DynamoTable *dynamoTable `json:"dynamo_table"`
	Cassandra   *cassandra   `json:"cassandra"`
	Jaeger      *jaeger      `json:"jaeger"`
}

func init() {
	viper.SetDefault("service", "tracker")
	viper.SetDefault("version", "0.1.0-alpha")
	viper.SetDefault("stage", DevelopmentStage)
}

func (c *Configuration) LoadEnv() error {
	//	rule
	//	if stage is development (dev) then use config file
	//	else use environment variables
	c.Stage = os.Getenv("LT_STAGE")
	if c.IsDevEnv() || c.IsTestEnv() {
		return c.loadFromFile()
	}

	viper.SetEnvPrefix("lt")
	viper.AutomaticEnv()
	return nil
}

func (c Configuration) IsDevEnv() bool {
	return c.Stage == DevelopmentStage
}

func (c Configuration) IsTestEnv() bool {
	return c.Stage == TestingStage
}

func (c Configuration) IsStagingEnv() bool {
	return c.Stage == StagingStage
}

func (c Configuration) IsProdEnv() bool {
	return c.Stage == ProductionStage
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
	c.Service = strings.ToLower(viper.GetString("service"))
	c.Version = strings.ToLower(viper.GetString("version"))
	c.Stage = c.getStage()
	c.HTTP.Load(c.Version)
	c.DynamoTable.Load(c.Stage)
	c.Cassandra.Load(c.Stage)
	c.Jaeger.Load()
}

func (c Configuration) getStage() string {
	stage := strings.ToLower(viper.GetString("stage"))
	isValidStage := stage == DevelopmentStage || stage == TestingStage || stage == StagingStage ||
		stage == ProductionStage
	if !isValidStage {
		return DevelopmentStage
	}

	return stage
}

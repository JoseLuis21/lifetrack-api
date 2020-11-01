package configuration

import "github.com/spf13/viper"

// dynamoTable AWS DynamoDB Table config
type dynamoTable struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

func init() {
	viper.SetDefault("dynamodb.table", "lifetrack-"+DevelopmentStage)
	viper.SetDefault("dynamodb.region", "us-east-1")
}

func (d *dynamoTable) Load(stage string) {
	d.Name = setResourceStage(viper.GetString("dynamodb.table"), stage)
	d.Region = viper.GetString("dynamodb.region")
}

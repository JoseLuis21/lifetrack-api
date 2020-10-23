package main

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("lifetrack.persistence.dynamo.table", "lifetrack-dev")
	viper.AutomaticEnv()
	if err := viper.BindEnv("lifetrack.persistence.dynamo.table", "LT_TABLE_NAME"); err != nil {
		panic(err)
	}
}

func main() {
	log.Print(viper.GetString("lifetrack.persistence.dynamo.table"))
}

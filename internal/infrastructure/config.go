package infrastructure

import "github.com/spf13/viper"

type Configuration struct {
	Table struct {
		Name   string
		Region string
	}
}

func NewConfiguration() (*Configuration, error) {
	viper.SetDefault("table_name", "lt-category-dev")
	viper.SetDefault("table_region", "us-east-1")

	if err := SetOSEnv(); err != nil {
		return nil, err
	}

	return &Configuration{
		Table: struct {
			Name   string
			Region string
		}{
			Name:   viper.GetString("table_name"),
			Region: viper.GetString("table_region"),
		},
	}, nil
}

func SetOSEnv() error {
	viper.SetEnvPrefix("lt")
	if err := viper.BindEnv("table_name"); err != nil {
		return err
	}

	return viper.BindEnv("table_region")
}

package infrastructure

import "github.com/spf13/viper"

type Configuration struct {
	Category struct {
		Table struct {
			Name   string
			Region string
		}
		Event struct {
			Region string
		}
	}
}

func NewConfiguration() (Configuration, error) {
	viper.SetDefault("category.table.name", "lt-category")
	viper.SetDefault("category.table.region", "us-east-1")
	viper.SetDefault("category.event.region", "us-east-1")

	if err := SetOSEnv(); err != nil {
		return Configuration{}, err
	}

	return Configuration{
		Category: struct {
			Table struct {
				Name   string
				Region string
			}
			Event struct {
				Region string
			}
		}{
			Table: struct {
				Name   string
				Region string
			}{
				Name:   viper.GetString("category.table.name"),
				Region: viper.GetString("category.table.region"),
			},
			Event: struct {
				Region string
			}{
				Region: viper.GetString("category.event.region"),
			},
		},
	}, nil
}

func SetOSEnv() error {
	viper.SetEnvPrefix("lt")
	if err := viper.BindEnv("category_table_name", "category.table.name"); err != nil {
		return err
	} else if err := viper.BindEnv("category_event_region", "category.event.region"); err != nil {
		return err
	}

	return viper.BindEnv("category_table_region", "category.table.region")
}

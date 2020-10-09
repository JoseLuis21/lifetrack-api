package infrastructure

import "github.com/spf13/viper"

type dynamo struct {
	Table  string
	Region string
}

type persistence struct {
	DynamoDB dynamo
}

type awsEvent struct {
	Region string
}

type eventBus struct {
	AWS awsEvent
}

type Configuration struct {
	Persistence persistence
	EventBus    eventBus
}

func NewConfiguration() (Configuration, error) {
	viper.SetDefault("lifetrack.persistence.dynamo.table", "lifetrack-prod")
	viper.SetDefault("lifetrack.persistence.dynamo.region", "us-east-1")
	viper.SetDefault("lifetrack.eventbus.aws.region", "us-east-1")

	if err := SetOSEnv(); err != nil {
		return Configuration{}, err
	}

	return Configuration{
		Persistence: persistence{DynamoDB: dynamo{
			Table:  viper.GetString("lifetrack.persistence.dynamo.table"),
			Region: viper.GetString("lifetrack.persistence.dynamo.region"),
		},
		},
		EventBus: eventBus{
			awsEvent{Region: viper.GetString("lifetrack.eventbus.aws.region")},
		},
	}, nil
}

func SetOSEnv() error {
	viper.SetEnvPrefix("lt")
	if err := viper.BindEnv("dynamo_table_name", "lifetrack.persistence.dynamo.table"); err != nil {
		return err
	} else if err := viper.BindEnv("event_aws_region", "lifetrack.eventbus.aws.region"); err != nil {
		return err
	}

	return viper.BindEnv("dynamo_table_region", "lifetrack.persistence.dynamo.region")
}

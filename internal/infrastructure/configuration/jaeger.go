package configuration

import "github.com/spf13/viper"

type jaeger struct {
	AgentEndpoint     string `json:"agent_endpoint"`
	CollectorEndpoint string `json:"collector_endpoint"`
}

func init() {
	viper.SetDefault("jaeger.agent", "localhost:6831")
	viper.SetDefault("jaeger.collector", "http://localhost:14268/api/traces")
}

func (j *jaeger) Load() {
	j.AgentEndpoint = viper.GetString("jaeger.agent")
	j.CollectorEndpoint = viper.GetString("jaeger.collector")
}

package observability

import (
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// NewHTTP handles observability depending on the current development stage, strategy pattern applied
func NewHTTP(logger *zap.Logger, cfg configuration.Configuration) {
	//	rules
	//	a.	if stage == dev or test, then use jaeger and prometheus OpenCensus exporters
	//	b.	if stage != dev or test, then use production config (AWS X-Ray and CloudWatch)
	switch {
	case cfg.IsDevEnv() || cfg.IsTestEnv():
		newJaegerHTTP(logger, cfg)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		// TODO: Add Prometheus
	default:
		newJaegerHTTP(logger, cfg)
		// TODO: Add AWS X-Ray
		// 1: 1,000 samples exported
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1 / 1000.0)})
		break
	}
}

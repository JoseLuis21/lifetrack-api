package observability

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// newJaegerHTTP register a Jaeger trace exporter to the current OpenCensus implementation
func newJaegerHTTP(logger *zap.Logger, cfg configuration.Configuration) {
	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     cfg.Jaeger.AgentEndpoint,
		CollectorEndpoint: cfg.Jaeger.CollectorEndpoint,
		Process: jaeger.Process{
			ServiceName: cfg.Service,
			Tags:        nil,
		},
	})
	if err != nil {
		logger.Warn("failed to initialize jaeger", zap.Error(err))
		return
	}
	trace.RegisterExporter(je)
	logger.Debug("registered jaeger exporter")
}

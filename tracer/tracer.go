package tracer

import (
	"context"
	"fmt"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/helper"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("gin-server")

func InitializeTracer() {
	tp, err := initTracer()
	if err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("%v", err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.GetLogger().Error(fmt.Sprintf("Error shutting down tracer provider: %v", err))
		}
	}()
}

func initTracer() (*sdktrace.TracerProvider, error) {
	jaegerExporter, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(helper.GetEnv("JAEGER_HOST", "localhost")), jaeger.WithAgentPort(helper.GetEnv("JAEGER_PORT", "14250"))))
	stdoutExporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(jaegerExporter),
		sdktrace.WithBatcher(stdoutExporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func GetTracer() trace.Tracer {
	return tracer
}

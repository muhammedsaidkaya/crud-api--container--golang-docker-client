package main

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/cache_layer"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/controller_layer"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/helper"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/logger"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/metric"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/repository_layer"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/service_layer"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/tracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	lemonsKey = attribute.Key("ex.com/lemons")
)

func main() {

	logger.InitializeLogger(helper.GetEnv("LOG_LEVEL", "INFO"), helper.GetEnv("LOG_FILE", "logfile"))
	enableTracingByOption()
	enableMetricByOption()

	router := gin.New()
	router.Use(otelgin.Middleware("gin-docker-client-api"))
	InitializeLayers().Setup(router)
	router.Run(":" + helper.GetEnv("APP_PORT", "8080"))
}

func InitializeLayers() controller_layer.ContainerControllerInterface {
	dockerClient, _ := client.NewEnvClient()
	cacheDefaultExpirationTime, _ := strconv.Atoi(helper.GetEnv("CACHE_DEFAULT_EXPIRATION_TIME", "10"))
	cacheCleanUpIntervalTime, _ := strconv.Atoi(helper.GetEnv("CACHE_CLEANUP_INTERVAL_TIME", "10"))
	repository := repository_layer.NewContainerRepository(cache_layer.NewCache(time.Duration(cacheDefaultExpirationTime)*time.Second, time.Duration(cacheCleanUpIntervalTime)*time.Second), dockerClient)
	service := service_layer.NewContainerService(repository)
	return controller_layer.NewContainerController(service)
}

func enableTracingByOption() {
	tracingEnabled := helper.GetEnv("JAEGER_EXPORTER_ENABLED", "false")
	boolValue, err := strconv.ParseBool(tracingEnabled)
	if err != nil {
		log.Fatal(err)
	}
	if boolValue {
		tracer.InitTracer()
	}
}

func enableMetricByOption() {
	metricEnabled := helper.GetEnv("PROMETHEUS_EXPORTER_ENABLED", "false")
	boolValue, err := strconv.ParseBool(metricEnabled)
	if err != nil {
		log.Fatal(err)
	}
	if boolValue {
		metric.InitMeter()
		createUpTimeMetric()
	}
}

func createUpTimeMetric() {
	meter := global.Meter("custom")

	observerLock := new(sync.RWMutex)
	observerValueToReport := new(float64)
	observerAttrsToReport := []attribute.KeyValue{attribute.String("app", "docker-client")}

	gaugeObserver, err := meter.AsyncFloat64().Gauge("uptime")
	if err != nil {
		log.Panicf("failed to initialize instrument: %v", err)
	}
	_ = meter.RegisterCallback([]instrument.Asynchronous{gaugeObserver}, func(ctx context.Context) {
		(*observerLock).RLock()
		value := *observerValueToReport
		attrs := observerAttrsToReport
		(*observerLock).RUnlock()
		gaugeObserver.Observe(ctx, value, attrs...)
	})

	go func() {
		for {
			*observerValueToReport = *observerValueToReport + 1.0
			time.Sleep(time.Second)
		}
	}()
}

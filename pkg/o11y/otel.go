package o11y

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/blackskad/quasar/pkg/environment"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	opts := []sdktrace.TracerProviderOption{
		// TODO: add ability to switch between sdktrace.AlwaysSample and sdktrace.ProbabilitySampler for production.
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(environment.Name),
			),
		),
	}

	// Always export to stdout
	expStdout, err := stdouttrace.New()
	if err != nil {
		return nil, err
	}
	opts = append(opts, sdktrace.WithBatcher(expStdout))

	// If the env var is set, also export to an otel collector
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		expOtel, err := otlptracegrpc.New(ctx)
		if err != nil {
			return nil, err
		}
		opts = append(opts, sdktrace.WithBatcher(expOtel))
	}

	// TODO: provider is not properly shutdown on program exit
	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

func initMeter() (*sdkmetric.MeterProvider, error) {
	// Translate all otel metrics into prometheus format
	exp, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	// TODO: provider is not properly shutdown on program exit
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exp))
	otel.SetMeterProvider(mp)

	// Expose the prometheus metrics for scraping
	go runPrometheusServer()

	// Record go runtime metrics like goroutines, heap & gc
	runtime.Start(runtime.WithMeterProvider(mp))

	return mp, nil
}

// runPProfServer starts an HTTP server on port 9090 that exposes the metrics on path /metrics.
// This function blocks until an error occurs
// This is a separate function to make stacktraces more readable
func runPrometheusServer() {
	mux := &http.ServeMux{}
	mux.Handle("/metrics", promhttp.Handler())

	// TODO: server is not properly shutdown on program exit
	err := http.ListenAndServe(":9090", mux)
	if err != nil {
		fmt.Printf("error serving prometheus: %v", err)
		return
	}
}

var initOnce sync.Once

func Register(ctx context.Context, h http.Handler) http.Handler {
	initOnce.Do(func() {
		_, err := initTracer(ctx)
		if err != nil {
			log.Fatal(err)
		}

		_, err = initMeter()
		if err != nil {
			log.Fatal(err)
		}
	})
	return otelhttp.NewHandler(h, "api")
}

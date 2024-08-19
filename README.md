# Go Web Template

## What's included

### Basic setup

The main entry point is in `/main.go`, however the web logic will live in `/pkg/web/` for endpoints that return data for human consumption in the browser (HTML, static images, CSS and JS,...) and `/pkg/web/api` for endpoints that return data for programatic consumption (JSON, XML, ...)

### Docker image

While the app can be build and run with straight Go, it's intended to be build and run inside a docker container. The Dockerfile contains four stages:

* Base: this layer copies all the files into the Go image and pulls in all the Go mod dependencies.

* Test: this layer builds off of the base layer and runs `go vet` and `go test`. 

* Build: this layer builds off of the base layer and builds the production binaries.

* Production: this layer builds off of the distroless image. It sets up the binary and opens up the required ports in the image.

### Docker compose

To allow you to easily run the service locally without much local config, a docker-compose.yaml file is included. This file will include everything to run a minimal stack.

### Profiling

The service always runs with pprof enabled on port 6060. This allows you to fetch runtime profiling information on `http://localhost:6060/debug/pprof`


### Observability

#### Metrics

The service records metrics through the metrics subsystem of Open Telemetry. It has an exporter configured that exposes all metrics on an endpoint `http://localhost:9090/metrics`, in a prometheus-compatible format.

The service will track the Go runtime metrics (goroutines, heap and GC stats) and metrics for the HTTP endpoints.

#### Tracing

The service has tracing set up for the main HTTP mux. All endpoints registered on this mux will create a new trace. During the handling of the endpoint, you can create a new span that is included on the trace of the endpoint by using 

```
	ctx, span := otel.GetTracerProvider().Tracer(environment.AppName).Start(r.Context(), "endpoint-sub-span")
	defer span.End()
```

The open telemetry traces are exported to stdout by default. If the Open Telemetry environment variable `OTEL_EXPORTER_OTLP_ENDPOINT` has been configured, a grpc-based OTLP exporter will be set up.
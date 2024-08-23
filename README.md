# Go Web Template

## What's included

### Basic setup

The main entry point is in `/main.go`, however the web logic will live in `/pkg/web/` for endpoints that return data for human consumption in the browser (HTML, static images, CSS and JS,...) and `/pkg/web/api` for endpoints that return data for programatic consumption (JSON, XML, ...)

### Docker image

While the app can be build and run with plain Go commands, it's intended to be build and run inside a docker container. The Dockerfile contains four stages:

* Base: this layer copies all the files into the Go image and pulls in all the Go mod dependencies.

* Test: this layer builds off of the base layer and runs `go vet` and `go test`. 

* Build: this layer builds off of the base layer and builds the production binaries.

* Production: this layer builds off of the distroless image. It sets up the binary and opens up the required ports in the image.

### Docker compose

To allow you to easily run the service locally without much local config, a docker-compose.yaml file is included. This file will include everything to run a minimal stack.

### CI pipeline - GitHub only

The project contains a GitHub Actions configuration file to run the docker build stages on `push` and `pull_requests`. It will first run the `test` stage, then it run the `production` stage.

While it will build the production image, the workflow is not configured to push the image to a docker image registry. You will have to uncomment the docker login job, change the `push` argument for the production job to `true`, and set a proper image tag.

In your GitHub repository, there is an option to configure rulesets for your main branch. Within the ruleset, the success of the build workflow can be made required for each pull requests. For more information, please see the [GitHub documentation](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets#require-status-checks-to-pass-before-merging).

### App configuration

The main configuration will be done through environment variables. The environment package will parse the environment variables into a struct that can then be passed around through the service.

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
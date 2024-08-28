# Go Web Scaffold

This Go Web Scaffold is a ready-to-go project that will get you going with your Go web application straight out of the door, with all the nuts and bolt that you expect in your projects: a decent project layout, tablestakes app plumbing, a local development setup, a proper CI pipeline with artifact publishing, and more...

The core philosophies of this scaffold are __simplicity__ and __efficiency__. The scaffold should be self-explanatory and easy to use for anyone who's seen a Go project before. You shouldn't need to search within the codebase for a piece of functionality. You shouldn't need to learn extra tools besides Go to get started. You shouldn't have to jump though a ton of configuration to get going. The scaffold should also not be doing anything that causes useless overhead, neither in the build process nor in the actual application. It shouldn't publish build artifacts that noone ever looks at. It shouldn't use CPU or memory hungry libraries by default.

## Getting started

How to get started with the scaffolding:

1. Clone this repository
	```
	go install golang.org/x/tools/cmd/gonew@latest
	gonew github.com/blackskad/go-web-scaffold github.com/$ACCOUNT/$PROJECTNAME
	```

1. Replace all mentions of go-web-scaffold with your own project name. At the moment, that's in the [Dockerfile](https://github.com/blackskad/go-web-scaffold/blob/main/Dockerfile#L25) and in the GitHub [build workflow](https://github.com/blackskad/go-web-scaffold/blob/main/.github/workflows/build.yml#L53).
	```
	sed -i 's/blackskad\/go-web-scaffold/$ACCOUNT\/$PROJECTNAME/g' .
	```

1. Initialize the git repo and push them to your own GitHub repo.
	```
	git init
	git remote add origin git@github.com:$ACCOUNT/$PROJECTNAME
	git push origin main
	```

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

To allow you to easily run the service locally without having to install a lot of local dependencies, a docker-compose.yaml file is included. This file includes everything to run a minimal stack.

Run `docker compose up --build` to make the service accessible [locally on port 8080](http://localhost:8080)

### GitHub CI pipeline

The project contains a GitHub Actions configuration file to run the docker build stages on `push` to any branch and semver tags, and on `pull_requests` against the `main` branch. It will first run the `test` stage, then it run the `production` stage.

When the pipeline runs for a semver tag, the tag will be embedded in the binary as the application version and the image will be pushed to ghcr. None of the other runs will push a docker image.

While it is recommended to follow [trunk-based development](https://trunkbaseddevelopment.com/), it is not enforced by the build system in any way.

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
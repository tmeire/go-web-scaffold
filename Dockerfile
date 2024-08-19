# Stage to copy the code and fetch dependencies 
# try to keep the debian version in sync with the distroless version
FROM golang:1.23-bookworm as base

ARG VERSION=development

WORKDIR /work

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main.go ./
COPY pkg/ ./pkg

# Stage to test the code
FROM base as test

RUN go vet -v
RUN go test -v

# Stage to build the binary
FROM base as build

RUN CGO_ENABLED=0 go build -ldflags="-X 'github.com/blackskad/go-web-scaffold/environment.Version=${VERSION}'" -o app .

# Stage with the production binary
FROM gcr.io/distroless/static-debian12 as production

COPY --from=build /work/app /

# port for pprof
EXPOSE 6060
# port for the actual app
EXPOSE 8080
# port for the prometheus metrics
EXPOSE 9090

CMD ["/app"]
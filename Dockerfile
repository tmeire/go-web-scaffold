# try to keep the debian version in sync with the distroless version
FROM golang:1.23-bookworm as base

WORKDIR /go/src/app

COPY go.mod ./
RUN go mod download

COPY main.go ./
COPY pkg/ ./pkg

FROM base as test

RUN go vet -v
RUN go test -v

FROM base as build

RUN CGO_ENABLED=0 go build -o /go/bin/app .

FROM gcr.io/distroless/static-debian12 as production

COPY --from=build /go/bin/app /
CMD ["/app"]
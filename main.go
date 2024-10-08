package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmeire/go-web-scaffold/pkg/environment"
	"github.com/tmeire/go-web-scaffold/pkg/o11y"
	"github.com/tmeire/go-web-scaffold/pkg/web"
)

func main() {
	ctx := context.Background()

	conf := environment.Parse()

	o11y.StartPProfServer()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	mux := http.NewServeMux()
	web.Register(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: o11y.Register(ctx, conf, mux),
	}

	done := make(chan struct{})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
		close(done)
	}()

	select {
	case <-sigc:
		err := server.Close()
		if err != nil {
			panic(err)
		}
	case <-done:
		return
	}
}

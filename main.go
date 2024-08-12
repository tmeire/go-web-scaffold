package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/blackskad/quasar/pkg/web"
)

func main() {
	go func() {
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			slog.Warn("pprof http server shut down", slog.Any("error", err))
		}
	}()

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
		Handler: mux,
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

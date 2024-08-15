package o11y

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
)

func StartPProfServer() {
	go runPProfServer()
}

// runPProfServer starts an HTTP server on port 6060 that exposes /debug/pprof
// This function blocks until an error occurs
// This is a separate function to make stacktraces more readable
func runPProfServer() {
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		slog.Warn("pprof http server shut down", slog.Any("error", err))
	}
}

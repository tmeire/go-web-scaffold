package o11y

import (
	_ "embed"
	"html/template"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
)

func StartPProfServer() {
	go runPProfServer()
}

//go:embed debug-info.html
var debugInfoPage string

var debugInfoTPL = template.Must(template.New("debug-info").Parse(debugInfoPage))

// runPProfServer starts an HTTP server on port 6060 that exposes /debug/pprof
// This function blocks until an error occurs
// This is a separate function to make stacktraces more readable
func runPProfServer() {
	http.HandleFunc("/debug/info", func(w http.ResponseWriter, r *http.Request) {
		debugInfo, ok := debug.ReadBuildInfo()
		if !ok {
			http.Error(w, "unable to load build info", http.StatusInternalServerError)
			return
		}

		err := debugInfoTPL.ExecuteTemplate(w, "debug-info", debugInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		slog.Warn("pprof http server shut down", slog.Any("error", err))
	}
}

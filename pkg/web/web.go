package web

import (
	"net/http"

	"github.com/tmeire/go-web-scaffold/pkg/web/auth"
)

func Register(h *http.ServeMux) {
	auth.Register(h)
}

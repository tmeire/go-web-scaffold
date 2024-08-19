package web

import (
	"net/http"

	"github.com/blackskad/go-web-scaffold/pkg/web/auth"
)

func Register(h *http.ServeMux) {
	auth.Register(h)
}

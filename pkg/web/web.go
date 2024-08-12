package web

import (
	"net/http"

	"github.com/blackskad/quasar/pkg/web/auth"
)

func Register(h *http.ServeMux) {
	auth.Register(h)
}

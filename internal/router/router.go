package router

import (
	"net/http"

	"threats-info-exchange/internal/handler"
)

func New(h *handler.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /info", h.PostInfo)
	mux.HandleFunc("GET /info/{md5}", h.GetInfo)
	mux.HandleFunc("PUT /info/{md5}", h.UpdateInfo)

	return mux
}

package api

import (
	"gps/pkg/ws"
	"net/http"
)

type handler struct {
	ws *ws.Manager
}

func NewHandler(wsManager *ws.Manager) *handler {
	return &handler{
		ws: wsManager,
	}
}

func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	// Sign up logic
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	// Login logic
}

func (h *handler) createRoute(w http.ResponseWriter, r *http.Request) {
	// Create route logic
}

func (h *handler) websocket(w http.ResponseWriter, r *http.Request) {
	// WebSocket handling logic
}

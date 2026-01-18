package api

import (
	"context"
	"gps/internal/adapters/api/middleware"
	"gps/pkg/ws"
	"net/http"
)

type Api struct {
	websocketManager *ws.Manager
	authMiddleware   middleware.Middleware
	handler          *handler
	server           *http.Server
}

func (a *Api) Start() error {
	mux := http.NewServeMux()

	mux.Handle("POST /sign_up", middleware.LoggingMiddleware(a.handler.signUp))
	mux.Handle("POST /login", middleware.LoggingMiddleware(a.handler.login))
	mux.Handle("GET /ws/{user_id}", middleware.LoggingMiddleware(a.handler.websocket))

	a.server.Handler = mux
	return a.server.ListenAndServe()
}

func (a *Api) StopServer(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *Api) GetWebsocketManager() *ws.Manager {
	return a.websocketManager
}

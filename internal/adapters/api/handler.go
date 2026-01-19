package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"gps/internal/app_services/auth"
	"gps/internal/domain/models"
	"gps/internal/domain/services"
	"gps/pkg/ws"

	"github.com/google/uuid"
)

type handler struct {
	ws         *ws.Manager
	auth       AuthService
	aggregator Aggregator
}

type AuthService interface {
	SignUp(ctx context.Context, input auth.Input) (string, error)
	LogIn(ctx context.Context, input auth.Input) (string, error)
}

type Aggregator interface {
	AggregateRoute(route models.Route) models.AggregatedData
}

func NewHandler(wsManager *ws.Manager, authService AuthService, aggregator Aggregator) *handler {
	if aggregator == nil {
		aggregator = services.NewAggregator()
	}
	return &handler{
		ws:         wsManager,
		auth:       authService,
		aggregator: aggregator,
	}
}

func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusNotImplemented, "auth service not configured")
		return
	}
	var input auth.Input
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.auth.SignUp(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusNotImplemented, "auth service not configured")
		return
	}
	var input auth.Input
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.auth.LogIn(r.Context(), input)
	if err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, auth.ErrUserNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// func (h *handler) createRoute(w http.ResponseWriter, r *http.Request) {

// 	var route models.Route
// 	if err := decodeJSON(r, &route); err != nil {
// 		writeError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if route.RouteID == uuid.Nil {
// 		route.RouteID = uuid.New()
// 	}
// 	if route.StartTime.IsZero() && len(route.Path) > 0 {
// 		route.StartTime = route.Path[0].Timestamp
// 	}
// 	if route.EndTime.IsZero() && len(route.Path) > 0 {
// 		route.EndTime = route.Path[len(route.Path)-1].Timestamp
// 	}

// 	if err := h.routes.StoreRoute(r.Context(), route); err != nil {
// 		writeError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	writeJSON(w, http.StatusCreated, map[string]string{"route_id": route.RouteID.String()})
// }

func (h *handler) websocket(w http.ResponseWriter, r *http.Request) {
	if h.ws == nil {
		writeError(w, http.StatusNotImplemented, "websocket manager not configured")
		return
	}
	userID, err := parseUUIDParam(r, "route_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.ws.ServeWS(w, r, userID)
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parseUUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	value := r.PathValue(name)
	if value == "" {
		return uuid.Nil, errors.New("missing path parameter")
	}
	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, err
	}
	return parsed, nil
}

package interfaces

import (
	"context"

	"gps/internal/domain/models"

	"github.com/google/uuid"
)

type RouteRepository interface {
	CreateRoute(ctx context.Context, route models.Route) error
	GetRouteByID(ctx context.Context, routeID uuid.UUID) (models.Route, error)
	AddGPSDataToRoute(ctx context.Context, routeID uuid.UUID, gps models.GPSData) error
	GetGPSDataLastNSeconds(ctx context.Context, seconds int) ([]models.GPSData, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, username, passwordHash string) (uuid.UUID, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}

type RedisRouteRepository interface {
	StoreRoute(ctx context.Context, route models.Route) error
	AppendRoutePoint(ctx context.Context, routeID uuid.UUID, point models.GPSData) error
	GetRoute(ctx context.Context, routeID uuid.UUID) (models.Route, error)
	GetRoutePath(ctx context.Context, routeID uuid.UUID) ([]models.GPSData, error)
	DeleteRoute(ctx context.Context, routeID uuid.UUID) error
}

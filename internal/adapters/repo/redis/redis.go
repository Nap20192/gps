package redisRepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gps/internal/domain/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	routeKeyPrefix  = "route:"
	routePathSuffix = ":gps"
	startTimeField  = "start_time"
	endTimeField    = "end_time"
)

type Repository struct {
	client *redis.Client
	ttl    time.Duration
}

var ErrRouteNotFound = errors.New("route not found")

func NewRepository(client *redis.Client, ttl time.Duration) (*Repository, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	return &Repository{client: client, ttl: ttl}, nil
}

func (r *Repository) StoreRoute(ctx context.Context, route models.Route) error {
	if r == nil || r.client == nil {
		return fmt.Errorf("redis repository is not initialized")
	}
	if route.RouteID == uuid.Nil {
		return fmt.Errorf("route id is required")
	}

	routeKey := routeMetaKey(route.RouteID)
	pathKey := routePathKey(route.RouteID)

	pipe := r.client.TxPipeline()
	pipe.Del(ctx, routeKey, pathKey)

	if !route.StartTime.IsZero() || !route.EndTime.IsZero() {
		fields := map[string]any{}
		if !route.StartTime.IsZero() {
			fields[startTimeField] = route.StartTime.UnixNano()
		}
		if !route.EndTime.IsZero() {
			fields[endTimeField] = route.EndTime.UnixNano()
		}
		pipe.HSet(ctx, routeKey, fields)
	}

	for _, point := range route.Path {
		payload, err := json.Marshal(point)
		if err != nil {
			pipe.Discard()
			return err
		}
		score := float64(point.Timestamp.UnixNano())
		pipe.ZAdd(ctx, pathKey, redis.Z{Score: score, Member: payload})
	}
	if r.ttl > 0 {
		pipe.Expire(ctx, routeKey, r.ttl)
		pipe.Expire(ctx, pathKey, r.ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Repository) AppendRoutePoint(ctx context.Context, routeID uuid.UUID, point models.GPSData) error {
	if r == nil || r.client == nil {
		return fmt.Errorf("redis repository is not initialized")
	}
	if routeID == uuid.Nil {
		return fmt.Errorf("route id is required")
	}

	routeKey := routeMetaKey(routeID)
	pathKey := routePathKey(routeID)
	payload, err := json.Marshal(point)
	if err != nil {
		return err
	}

	score := float64(point.Timestamp.UnixNano())
	pipe := r.client.TxPipeline()
	pipe.ZAdd(ctx, pathKey, redis.Z{Score: score, Member: payload})
	pipe.HSet(ctx, routeKey, endTimeField, point.Timestamp.UnixNano())
	pipe.HSetNX(ctx, routeKey, startTimeField, point.Timestamp.UnixNano())

	if r.ttl > 0 {
		pipe.Expire(ctx, routeKey, r.ttl)
		pipe.Expire(ctx, pathKey, r.ttl)
	}

	_, err = pipe.Exec(ctx)
	return err
}

func (r *Repository) GetRoute(ctx context.Context, routeID uuid.UUID) (models.Route, error) {
	if r == nil || r.client == nil {
		return models.Route{}, fmt.Errorf("redis repository is not initialized")
	}
	if routeID == uuid.Nil {
		return models.Route{}, fmt.Errorf("route id is required")
	}

	pathKey := routePathKey(routeID)
	items, err := r.client.ZRange(ctx, pathKey, 0, -1).Result()
	if err != nil {
		return models.Route{}, err
	}
	if len(items) == 0 {
		return models.Route{}, ErrRouteNotFound
	}

	path := make([]models.GPSData, 0, len(items))
	for _, raw := range items {
		var point models.GPSData
		if err := json.Unmarshal([]byte(raw), &point); err != nil {
			return models.Route{}, err
		}
		path = append(path, point)
	}

	start, end := r.resolveRouteBounds(ctx, routeID, path)

	return models.Route{
		RouteID:   routeID,
		StartTime: start,
		EndTime:   end,
		Path:      path,
	}, nil
}

func (r *Repository) GetRoutePath(ctx context.Context, routeID uuid.UUID) ([]models.GPSData, error) {
	if r == nil || r.client == nil {
		return nil, fmt.Errorf("redis repository is not initialized")
	}
	if routeID == uuid.Nil {
		return nil, fmt.Errorf("route id is required")
	}

	pathKey := routePathKey(routeID)
	items, err := r.client.ZRange(ctx, pathKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	path := make([]models.GPSData, 0, len(items))
	for _, raw := range items {
		var point models.GPSData
		if err := json.Unmarshal([]byte(raw), &point); err != nil {
			return nil, err
		}
		path = append(path, point)
	}
	return path, nil
}

func (r *Repository) DeleteRoute(ctx context.Context, routeID uuid.UUID) error {
	if r == nil || r.client == nil {
		return fmt.Errorf("redis repository is not initialized")
	}
	if routeID == uuid.Nil {
		return fmt.Errorf("route id is required")
	}

	routeKey := routeMetaKey(routeID)
	pathKey := routePathKey(routeID)
	return r.client.Del(ctx, routeKey, pathKey).Err()
}

func routeMetaKey(routeID uuid.UUID) string {
	return routeKeyPrefix + routeID.String()
}

func routePathKey(routeID uuid.UUID) string {
	return routeKeyPrefix + routeID.String() + routePathSuffix
}

func (r *Repository) resolveRouteBounds(ctx context.Context, routeID uuid.UUID, path []models.GPSData) (time.Time, time.Time) {
	routeKey := routeMetaKey(routeID)
	fields, err := r.client.HMGet(ctx, routeKey, startTimeField, endTimeField).Result()
	if err == nil && len(fields) == 2 {
		start := parseUnixNano(fields[0])
		end := parseUnixNano(fields[1])
		if !start.IsZero() || !end.IsZero() {
			return start, end
		}
	}
	if len(path) == 0 {
		return time.Time{}, time.Time{}
	}
	return path[0].Timestamp, path[len(path)-1].Timestamp
}

func parseUnixNano(value any) time.Time {
	switch v := value.(type) {
	case string:
		if v == "" {
			return time.Time{}
		}
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}
		}
		return time.Unix(0, parsed)
	case int64:
		return time.Unix(0, v)
	case float64:
		return time.Unix(0, int64(v))
	default:
		return time.Time{}
	}
}

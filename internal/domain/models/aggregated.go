package models

import (
	"time"

	"github.com/google/uuid"
)

type AggregatedData struct {
	RouteID      uuid.UUID     `json:"route_id" bson:"route_id"`
	AverageSpeed  float64       `json:"average_speed" bson:"average_speed"`
	TotalDistance float64       `json:"total_distance" bson:"total_distance"`
	Duration      time.Duration `json:"duration" bson:"duration"`
	AmountPoints  int           `json:"amount_points" bson:"amount_points"`
	Timestamp     time.Time     `json:"timestamp" bson:"timestamp"`
}

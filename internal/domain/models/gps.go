package models

import (
	"time"

	"github.com/google/uuid"
)

type GPSData struct {
	Location  Location  `json:"location" bson:"location"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Altitude  float64 `json:"altitude" bson:"altitude"`
}

type Route struct {
	RouteID   uuid.UUID `json:"route_id" bson:"route_id"`
	Path      []GPSData `json:"path" bson:"path"`
	StartTime time.Time `json:"start_time" bson:"start_time"`
	Finished  bool      `json:"finished" bson:"finished"`
	EndTime   time.Time `json:"end_time" bson:"end_time,omitempty"`
}

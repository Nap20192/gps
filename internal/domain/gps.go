package domain

import (
	"time"

	"github.com/google/uuid"
)

type GPSData struct {
	EntityID  uuid.UUID
	Location  Location
	Timestamp time.Time
}

type Location struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

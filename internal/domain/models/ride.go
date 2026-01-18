package models

import "github.com/google/uuid"

type Route struct {
	ID        uuid.UUID
	DriverID  uuid.UUID
	PassengerID uuid.UUID
	StartLocation Location
	EndLocation   Location
}

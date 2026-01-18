package interfaces

import (
	"context"
	"gps/internal/domain/models"
)

type Repository interface {
	StoreGPSData(ctx context.Context, data models.GPSData) error
	StoreAggregatedData(ctx context.Context, data models.AggregatedData) error

	GetGPSDataByUserID(ctx context.Context, userID string) ([]models.GPSData, error)
	GetAggregatedDataByUserID(ctx context.Context, userID string) ([]models.AggregatedData, error)
}

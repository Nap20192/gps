package interfaces

import (
	"context"
	"gps/internal/domain"
)

type Repository interface {
	StoreGPSData(ctx context.Context, data domain.GPSData) error
	StoreAggregatedData(ctx context.Context, data domain.AggregatedData) error

	GetGPSDataByUserID(ctx context.Context, userID string) ([]domain.GPSData, error)
	GetAggregatedDataByUserID(ctx context.Context, userID string) ([]domain.AggregatedData, error)

}

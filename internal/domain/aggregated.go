package domain

import "time"

type AggregatedData struct {
	AverageSpeed  float64
	TotalDistance float64
	Duration      time.Duration
	CreatedAt     time.Time
}

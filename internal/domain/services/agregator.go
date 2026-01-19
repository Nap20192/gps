package services

import (
	"math"
	"time"

	"gps/internal/domain/models"
)

const earthRadiusMeters = 6371000.0

type Aggregator struct{}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) AggregateRoute(route models.Route) models.AggregatedData {
	points := route.Path
	amountPoints := len(points)
	if amountPoints == 0 {
		return models.AggregatedData{
			RouteID:       route.RouteID,
			Timestamp:     a.resolveTimestamp(route, time.Time{}),
		}
	}

	totalDistance := 0.0
	for i := 1; i < amountPoints; i++ {
		prev := points[i-1]
		curr := points[i]
		totalDistance += distanceMeters(prev.Location, curr.Location)
	}

	startTime, endTime := a.resolveDurationBounds(route, points)
	duration := endTime.Sub(startTime)
	if duration < 0 {
		duration = 0
	}

	avgSpeed := 0.0
	seconds := duration.Seconds()
	if seconds > 0 {
		avgSpeed = totalDistance / seconds
	}

	return models.AggregatedData{
		RouteID:       route.RouteID,
		AverageSpeed:  avgSpeed,
		TotalDistance: totalDistance,
		Duration:      duration,
		AmountPoints:  amountPoints,
		Timestamp:     a.resolveTimestamp(route, points[amountPoints-1].Timestamp),
	}
}

func (a *Aggregator) resolveDurationBounds(route models.Route, points []models.GPSData) (time.Time, time.Time) {
	start := route.StartTime
	end := route.EndTime
	if start.IsZero() && len(points) > 0 {
		start = points[0].Timestamp
	}
	if end.IsZero() && len(points) > 0 {
		end = points[len(points)-1].Timestamp
	}
	if end.IsZero() {
		end = time.Now()
	}
	if start.IsZero() {
		start = end
	}
	return start, end
}

func (a *Aggregator) resolveTimestamp(route models.Route, fallback time.Time) time.Time {
	if !route.EndTime.IsZero() {
		return route.EndTime
	}
	if !fallback.IsZero() {
		return fallback
	}
	if !route.StartTime.IsZero() {
		return route.StartTime
	}
	return time.Now()
}

func distanceMeters(aLoc, bLoc models.Location) float64 {
	lat1 := toRadians(aLoc.Latitude)
	lat2 := toRadians(bLoc.Latitude)
	dLat := lat2 - lat1
	dLon := toRadians(bLoc.Longitude - aLoc.Longitude)

	sinLat := math.Sin(dLat / 2)
	sinLon := math.Sin(dLon / 2)

	h := sinLat*sinLat + math.Cos(lat1)*math.Cos(lat2)*sinLon*sinLon
	centralAngle := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
	horizontal := earthRadiusMeters * centralAngle

	altDelta := bLoc.Altitude - aLoc.Altitude
	if altDelta == 0 {
		return horizontal
	}
	return math.Sqrt(horizontal*horizontal + altDelta*altDelta)
}

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

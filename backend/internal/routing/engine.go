package routing

import (
	"math"
	"wikiwalkme-backend/internal/wikidata"
)

const EarthRadiusKm = 6371.0

// RouteResponse represents the final package we send to the frontend
type RouteResponse struct {
	TotalDistanceKm float64           `json:"total_distance_km"`
	EstimatedMins   float64           `json:"estimated_minutes"`
	Path            []wikidata.Target `json:"path"`
}

func GenerateRoute(startLat, startLon, endLat, endLon float64, maxMinutes float64, targets []wikidata.Target) RouteResponse {
	// Assume average walking speed of 5 km/h
	maxDistanceKm := (maxMinutes / 60.0) * 5.0

	var optimizedPath []wikidata.Target
	unvisited := make([]wikidata.Target, len(targets))
	copy(unvisited, targets)

	currentLat := startLat
	currentLon := startLon
	totalDistance := 0.0

	for len(unvisited) > 0 {
		bestIndex := -1
		minDistToTarget := math.MaxFloat64

		// 1. Find closest unvisited target to current location
		for i, target := range unvisited {
			dist := CalculateDistance(currentLat, currentLon, target.Lat, target.Lon)
			if dist < minDistToTarget {
				minDistToTarget = dist
				bestIndex = i
			}
		}

		if bestIndex == -1 {
			break
		}

		proposedTarget := unvisited[bestIndex]

		// 2. Calculate the detour cost: Current -> Proposed Target -> Final Destination
		distanceToTarget := CalculateDistance(currentLat, currentLon, proposedTarget.Lat, proposedTarget.Lon)
		distanceFromTargetToEnd := CalculateDistance(proposedTarget.Lat, proposedTarget.Lon, endLat, endLon)

		// 3. Check if this fits in our budget
		if totalDistance+distanceToTarget+distanceFromTargetToEnd <= maxDistanceKm {
			// It fits, commit to this step
			totalDistance += distanceToTarget
			optimizedPath = append(optimizedPath, proposedTarget)

			// Move current position to this target
			currentLat = proposedTarget.Lat
			currentLon = proposedTarget.Lon

			// Remove from unvisited pool
			unvisited = append(unvisited[:bestIndex], unvisited[bestIndex+1:]...)
		} else {
			// If closest point pushes us over out budget, we are done
			break
		}
	}

	// 4. Finally, connect the last point to the End Point
	totalDistance += CalculateDistance(currentLat, currentLon, endLat, endLon)
	estimatedMins := (totalDistance / 5.0) * 60.0

	return RouteResponse{
		TotalDistanceKm: totalDistance,
		EstimatedMins:   estimatedMins,
		Path:            optimizedPath,
	}
}

// Use Haversine formula to find distance between two points in kilometers
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Haversine formula math
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c
}

package routing

import (
	"math"
	"wikiwalkme-backend/internal/wikidata"
)

// RouteResponse represents the final package we send to the frontend
type RouteResponse struct {
	TotalDistanceKm float64           `json:"total_distance_km"`
	EstimatedMins   float64           `json:"estimated_minutes"`
	Path            []wikidata.Target `json:"path"`
}

func DistanceFlatEarth(lat1, lon1, lat2, lon2 float64) float64 {
	const degtoKm = 111.32
	dx := (lon2 - lon1) * math.Cos((lat1+lat2)*math.Pi/360.0) * degtoKm
	dy := (lat2 - lat1) * degtoKm
	return math.Sqrt(dx*dx + dy*dy)
}

func GenerateRoute(startLat, startLon, endLat, endLon float64, maxMinutes float64, allTargets []wikidata.Target) RouteResponse {
	// Convert minutes to max walking kilometers (assuing average walking pace of 4.5 km/h)
	maxKm := (maxMinutes / 60.0) * 4.5
	remainingKm := maxKm

	// Initial pruning: only keep targets reachable from start that still allow reaching end
	var reachableTargets []wikidata.Target
	for _, t := range allTargets {
		distFromStart := DistanceFlatEarth(startLat, startLon, t.Lat, t.Lon)
		distToEnd := DistanceFlatEarth(t.Lat, t.Lon, endLat, endLon)

		if distFromStart+distToEnd <= maxKm {
			reachableTargets = append(reachableTargets, t)
		}
	}

	// Greedy path construction - keep adding targets until none fit in remaining budget
	var path []wikidata.Target
	visited := make(map[string]bool)
	currLat, currLon := startLat, startLon

	for {
		var nextTarget *wikidata.Target
		bestDist := math.MaxFloat64

		// Find the closest unvisited target from pruned list that fits in remaining budget
		for i := range reachableTargets {
			t := &reachableTargets[i]
			if visited[t.WikidataUrl] {
				continue
			}

			distFromCurr := DistanceFlatEarth(currLat, currLon, t.Lat, t.Lon)
			distToEnd := DistanceFlatEarth(t.Lat, t.Lon, endLat, endLon)

			// Check if we can reach this target and still get to the end within remaining budget
			if distFromCurr+distToEnd <= remainingKm && distFromCurr < bestDist {
				bestDist = distFromCurr
				nextTarget = t
			}
		}

		// If no target found, we're done
		if nextTarget == nil {
			break
		}

		// Add target to path
		visited[nextTarget.WikidataUrl] = true
		path = append(path, *nextTarget)

		// Update position and remaining budget
		distTraveled := DistanceFlatEarth(currLat, currLon, nextTarget.Lat, nextTarget.Lon)
		remainingKm -= distTraveled
		currLat = nextTarget.Lat
		currLon = nextTarget.Lon
	}

	return RouteResponse{
		TotalDistanceKm: maxKm,
		EstimatedMins:   maxMinutes,
		Path:            path,
	}
}

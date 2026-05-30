package routing

import (
	"fmt"
	"math"
	"sort"
	"wikiwalkme-backend/internal/wikidata"
)

// GeoPoint is a simple coordinate for route drawing.
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// RouteResponse represents the final package we send to the frontend.
type RouteResponse struct {
	AllowedDistanceKm float64           `json:"allowed_distance_km"`
	TravelDistanceKm  float64           `json:"travel_distance_km"`
	EstimatedMins     float64           `json:"estimated_minutes"`
	Path              []wikidata.Target `json:"path"`
	RoutePoints       []GeoPoint        `json:"route_points"`
}

func DistanceFlatEarth(lat1, lon1, lat2, lon2 float64) float64 {
	const degtoKm = 111.32
	dx := (lon2 - lon1) * math.Cos((lat1+lat2)*math.Pi/360.0) * degtoKm
	dy := (lat2 - lat1) * degtoKm
	return math.Sqrt(dx*dx + dy*dy)
}

func GenerateRoute(startLat, startLon, endLat, endLon float64, maxMinutes float64, allTargets []wikidata.Target) RouteResponse {
	allowedKm := (maxMinutes / 60.0) * 4.5

	// Initial pruning: remove unreachable targets that cannot be visited on the way to the end.
	var reachableTargets []wikidata.Target
	for _, t := range allTargets {
		distFromStart := DistanceFlatEarth(startLat, startLon, t.Lat, t.Lon)
		distToEnd := DistanceFlatEarth(t.Lat, t.Lon, endLat, endLon)
		if distFromStart+distToEnd <= allowedKm {
			reachableTargets = append(reachableTargets, t)
		}
	}

	// Sort candidate targets by distance from the start to keep the route focused.
	sort.Slice(reachableTargets, func(i, j int) bool {
		iDist := DistanceFlatEarth(startLat, startLon, reachableTargets[i].Lat, reachableTargets[i].Lon)
		jDist := DistanceFlatEarth(startLat, startLon, reachableTargets[j].Lat, reachableTargets[j].Lon)
		return iDist < jDist
	})

	var path []wikidata.Target
	visited := make(map[string]bool)
	currLat, currLon := startLat, startLon
	travelKm := 0.0

	for {
		var nextTarget *wikidata.Target
		bestDist := math.MaxFloat64

		for i := range reachableTargets {
			t := &reachableTargets[i]
			if visited[t.WikidataUrl] {
				continue
			}

			distFromCurr := DistanceFlatEarth(currLat, currLon, t.Lat, t.Lon)
			distToEnd := DistanceFlatEarth(t.Lat, t.Lon, endLat, endLon)
			if distFromCurr+distToEnd <= allowedKm-travelKm && distFromCurr < bestDist {
				bestDist = distFromCurr
				nextTarget = t
			}
		}

		if nextTarget == nil {
			break
		}

		visited[nextTarget.WikidataUrl] = true
		path = append(path, *nextTarget)
		distanceToTarget := DistanceFlatEarth(currLat, currLon, nextTarget.Lat, nextTarget.Lon)
		travelKm += distanceToTarget
		currLat = nextTarget.Lat
		currLon = nextTarget.Lon
	}

	routePoints, err := FetchWalkingPath(startLat, startLon, endLat, endLon, path)
	if err != nil {
		fmt.Printf("OSRM fallback triggered: %v\n", err)
		routePoints = buildRoutePoints(startLat, startLon, endLat, endLon, path)
	}

	return RouteResponse{
		AllowedDistanceKm: allowedKm,
		TravelDistanceKm:  travelKm + DistanceFlatEarth(currLat, currLon, endLat, endLon),
		EstimatedMins:     (travelKm + DistanceFlatEarth(currLat, currLon, endLat, endLon)) / 4.5 * 60,
		Path:              path,
		RoutePoints:       routePoints,
	}
}

func buildRoutePoints(startLat, startLon, endLat, endLon float64, path []wikidata.Target) []GeoPoint {
	points := make([]GeoPoint, 0, len(path)+2)
	points = append(points, GeoPoint{Lat: startLat, Lon: startLon})
	for _, target := range path {
		points = append(points, GeoPoint{Lat: target.Lat, Lon: target.Lon})
	}
	points = append(points, GeoPoint{Lat: endLat, Lon: endLon})
	return points
}

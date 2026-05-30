package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"wikiwalkme-backend/internal/wikidata"
)

type OSRMResponse struct {
	Routes []struct {
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func FetchWalkingPath(startLat, startLon, endLat, endLon float64, path []wikidata.Target) ([]GeoPoint, error) {
	var coordPairs []string

	coordPairs = append(coordPairs, fmt.Sprintf("%f,%f", startLon, startLat))

	for _, t := range path {
		coordPairs = append(coordPairs, fmt.Sprintf("%f,%f", t.Lon, t.Lat))
	}

	coordPairs = append(coordPairs, fmt.Sprintf("%f,%f", endLon, endLat))

	coordString := strings.Join(coordPairs, ";")
	osrmURL := fmt.Sprintf("https://router.project-osrm.org/route/v1/walking/%s?overview=full&geometries=geojson", coordString)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(osrmURL)
	if err != nil {
		return nil, fmt.Errorf("failed to contact OSRM: %w", err)
	}
	defer resp.Body.Close()

	var result OSRMResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode OSRM response: %w", err)
	}

	if len(result.Routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	var geoPoints []GeoPoint
	for _, coord := range result.Routes[0].Geometry.Coordinates {
		if len(coord) == 2 {
			geoPoints = append(geoPoints, GeoPoint{Lat: coord[1], Lon: coord[0]})
		}
	}

	return geoPoints, nil
}

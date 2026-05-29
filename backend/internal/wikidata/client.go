package wikidata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Target struct {
	Name        string  `json:"name"`
	WikidataUrl string  `json:"wikidata_url"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

type WikidataResponse struct {
	Results struct {
		Bindings []struct {
			Item      struct{ Value string } `json:"item"`
			ItemLabel struct{ Value string } `json:"itemLabel"`
			Coords    struct{ Value string } `json:"coords"`
		} `json:"bindings"`
	} `json:"results"`
}

type WikidataClient struct {
	Cache *MemoryCache
}

func NewClient(cacheTTL time.Duration) *WikidataClient {
	return &WikidataClient{
		Cache: NewMemoryCache(cacheTTL),
	}
}

func (wc *WikidataClient) FetchTargets(lat, lon, radius float64) ([]Target, error) {
	// Check cache
	if cachedTargets, found := wc.Cache.Get(lat, lon); found {
		return cachedTargets, nil
	}

	sparqURL := "https://query.wikidata.org/sparql"
	query := fmt.Sprintf(`
		SELECT ?item ?itemLabel ?coords ?distance WHERE {
			SERVICE wikibase:around {
				?item wdt:P625 ?coords .
				bd:serviceParam wikibase:center "Point(%f %f)"^^geo:wktLiteral .
				bd:serviceParam wikibase:radius "%f" .
				bd:serviceParam wikibase:distance ?distance . # 1. Bind the calculated distance to ?distance
			}
			MINUS { ?item wdt:P18 ?image . }
			SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],en". }
		}
		ORDER BY ASC(?distance) # 2. Sort from closest to furthest
		LIMIT 200`,
		lon, lat, radius)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", sparqURL, nil)
	if err != nil {
		return nil, err
	}

	email := os.Getenv("WIKIDATA_CONTACT_EMAIL")
	if email == "" {
		email = "example@example.com"
	}
	useragent := fmt.Sprintf("WikiWalkMe/1.0 (%s) Go/net/http", email)
	req.Header.Set("User-Agent", useragent)
	req.Header.Set("Accept", "application/json")

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wikiResp WikidataResponse
	if err := json.NewDecoder(resp.Body).Decode(&wikiResp); err != nil {
		return nil, err
	}

	var targets []Target
	for _, binding := range wikiResp.Results.Bindings {
		cleaned := strings.Replace(binding.Coords.Value, "Point(", "", 1)
		cleaned = strings.Replace(cleaned, ")", "", 1)

		var itemLon, itemLat float64
		if _, err := fmt.Sscanf(cleaned, "%f %f", &itemLon, &itemLat); err == nil {
			targets = append(targets, Target{
				Name:        binding.ItemLabel.Value,
				WikidataUrl: binding.Item.Value,
				Lat:         itemLat,
				Lon:         itemLon,
			})
		}
	}

	// Save to cache
	wc.Cache.Set(lat, lon, targets)

	return targets, nil
}

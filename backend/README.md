# WikiWalkMe Backend

The backend engine is responsible for data discovery, route optimization, and fetching actual walking path geometry.

## Setup

1. Install Go 1.20+ from [https://go.dev/](https://go.dev/).
2. Run the server:
   ```bash
    cd backend
    go run ./cmd/server/main.go
   ```

## Key Features

- **Wikidata Integration**: Uses SPARQL queries to locate nearby points of interest.
- **Route Optimization**: Implements a distance-based pruning algorithm to fit stops into a user's time constraints.
- **Street-Level Routing**: Integrates with the OSRM (OpenStreetMap) API to generate walkable path geometry.
- **Memory Caching**: Implements a thread-safe `MemoryCache` to reduce redundant calls to the Wikidata SPARQL endpoint and OSRM.

## API Endpoints

- `POST /api/targets`: Returns a list of nearby photo targets based on coordinates.
- `POST /api/route`: Calculates an optimized route and returns street-level geometry.

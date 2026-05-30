# WikiWalkMe

Doesn't it just suck that there's so many places without photos on Wikimedia? Don't you also want to go on a little walk?

WikiWalkMe is a mapping application—inspired by [WikiShootMe!](https://wikishootme.toolforge.org/)—that bridges the gap between Wikidata discovery and actual pedestrian routing. It identifies nearby locations in need of documentation and calculates a walking path optimized to help you visit them as many of them as you can.

## Core Features

- **Target Discovery**: Queries Wikidata in real-time to find landmarks near your location lacking photographic documentation.
- **Path Optimization**: Uses a custom Go-based engine to prune and sort stops based on your available walking time.
- **Street-Level Routing**: Integrates with the OSRM (OpenStreetMap) walking profile to generate accurate, turn-by-turn path geometry.
- **Performance First**: Implements an in-memory caching layer for both Wikidata queries and route calculations, minimizing external API latency.

## Architecture

The system is divided into two primary modules:

### 1. Backend (`/backend`)

A Go-based API server powered by [Gin](https://github.com/gin-gonic/gin).

- **Engines**: Contains the core logic for spatial math and route pruning.
- **Clients**: Manages external connections to Wikidata and OSRM.
- **Caching**: A thread-safe memory cache that rounds coordinates to ensure frequent requests for the same area are served instantly.

### 2. Frontend (`/frontend`)

A React application using [Leaflet](https://leafletjs.com/) for interactive mapping.

- **State Management**: Uses custom hooks (`useTargets`, `useWalkingRoute`) to clean up component logic and handle asynchronous data fetching.
- **Interactive UI**: Provides a real-time slider for time management and simple click-to-route functionality.

## Intended Features (Roadmap)

We have a vision for making WikiWalkMe. Future iterations will include:

- **Persistent User Progress**: Save your completed routes and photo uploads so you can pick up where you left off.
- **Photo Upload Integration**: Directly upload photos to Wikimedia Commons from the app, linking them to the corresponding Wikidata items.
- **Advanced Target Filtering**: Toggle targets by type (e.g., historical buildings vs. natural monuments) or by specific maintenance categories.

## Getting Started

### Prerequisites

- [Go 1.20+](https://go.dev/)
- [Node.js 18+](https://nodejs.org/)

### Quick Start

1. **Backend**:
   ```bash
   cd backend
   go run main.go
   ```
2. **Frontend**:
   ```bash
    cd frontend
    npm install
    npm start
   ```

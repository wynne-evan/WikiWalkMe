package api

import (
	"net/http"
	"sort"
	"wikiwalkme-backend/internal/routing"
	"wikiwalkme-backend/internal/wikidata"

	"github.com/gin-gonic/gin"
)

type TargetRequest struct {
	Lat    float64 `json:"lat" binding:"required"`
	Lon    float64 `json:"lon" binding:"required"`
	Radius float64 `json:"radius"`
}

type RouteRequest struct {
	StartLat     float64 `json:"start_lat" binding:"required"`
	StartLon     float64 `json:"start_lon" binding:"required"`
	EndLat       float64 `json:"end_lat" binding:"required"`
	EndLon       float64 `json:"end_lon" binding:"required"`
	MaxMinutes   float64 `json:"max_minutes" binding:"required"`
	TargetRadius float64 `json:"target_radius"`
	MaxTargets   int     `json:"max_targets"`
}

type APIContext struct {
	WikiClient *wikidata.WikidataClient
	RouteCache *routing.RouteCache
}

func (api *APIContext) GenerateRouteHandler(c *gin.Context) {
	var req RouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cachedRoute, found := api.RouteCache.Get(req.StartLat, req.StartLon, req.EndLat, req.EndLon, req.MaxMinutes); found {
		c.JSON(http.StatusOK, cachedRoute)
		return
	}

	if req.TargetRadius <= 0 {
		req.TargetRadius = 5.0
	}
	if req.TargetRadius > 10.0 {
		req.TargetRadius = 10.0
	}

	startTargets, err := api.WikiClient.FetchTargets(req.StartLat, req.StartLon, req.TargetRadius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed fetching wikidata for start"})
		return
	}

	endTargets, err := api.WikiClient.FetchTargets(req.EndLat, req.EndLon, req.TargetRadius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed fetching wikidata for end"})
		return
	}

	allTargets := mergeTargets(startTargets, endTargets)

	if req.MaxTargets > 0 && len(allTargets) > req.MaxTargets {
		sort.Slice(allTargets, func(i, j int) bool {
			iDist := routing.DistanceFlatEarth(req.StartLat, req.StartLon, allTargets[i].Lat, allTargets[i].Lon)
			jDist := routing.DistanceFlatEarth(req.StartLat, req.StartLon, allTargets[j].Lat, allTargets[j].Lon)
			return iDist < jDist
		})
		allTargets = allTargets[:req.MaxTargets]
	}

	routeResult := routing.GenerateRoute(req.StartLat, req.StartLon, req.EndLat, req.EndLon, req.MaxMinutes, allTargets)
	api.RouteCache.Set(req.StartLat, req.StartLon, req.EndLat, req.EndLon, req.MaxMinutes, routeResult)
	c.JSON(http.StatusOK, routeResult)
}

func (api *APIContext) GetTargetsHandler(c *gin.Context) {
	var req TargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'lat' or 'lon'"})
		return
	}

	if req.Radius <= 0 {
		req.Radius = 10.0
	}
	if req.Radius > 10.0 {
		req.Radius = 10.0
	}

	targets, err := api.WikiClient.FetchTargets(req.Lat, req.Lon, req.Radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   len(targets),
		"targets": targets,
	})
}

func mergeTargets(startTargets, endTargets []wikidata.Target) []wikidata.Target {
	targetByURL := make(map[string]wikidata.Target)
	for _, target := range append(startTargets, endTargets...) {
		if target.WikidataUrl == "" {
			continue
		}
		targetByURL[target.WikidataUrl] = target
	}
	merged := make([]wikidata.Target, 0, len(targetByURL))
	for _, target := range targetByURL {
		merged = append(merged, target)
	}
	return merged
}

package api

import (
	"net/http"
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
	StartLat   float64 `json:"start_lat" binding:"required"`
	StartLon   float64 `json:"start_lon" binding:"required"`
	EndLat     float64 `json:"end_lat" binding:"required"`
	EndLon     float64 `json:"end_lon" binding:"required"`
	MaxMinutes float64 `json:"max_minutes" binding:"required"`
}

type APIContext struct {
	WikiClient *wikidata.WikidataClient
}

func (api *APIContext) GenerateRouteHandler(c *gin.Context) {
	var req RouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Fetch points around the starting area
	targets, err := api.WikiClient.FetchTargets(req.StartLat, req.StartLon, 3.0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed fetching wikidata"})
		return
	}

	// 2. Run the greedy path optimization
	routeResult := routing.GenerateRoute(req.StartLat, req.StartLon, req.EndLat, req.EndLon, req.MaxMinutes, targets)

	c.JSON(http.StatusOK, routeResult)
}

func (api *APIContext) GetTargetsHandler(c *gin.Context) {
	var req TargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'lat' or 'lon'"})
		return
	}

	if req.Radius == 0 {
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

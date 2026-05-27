package api

import (
	"net/http"
	"wikiwalkme-backend/internal/routing"
	"wikiwalkme-backend/internal/wikidata"

	"github.com/gin-gonic/gin"
)

type RouteRequest struct {
	StartLat   float64 `json:"start_lat" binding:"required"`
	StartLon   float64 `json:"start_lon" binding:"required"`
	EndLat     float64 `json:"end_lat" binding:"required"`
	EndLon     float64 `json:"end_lon" binding:"required"`
	MaxMinutes float64 `json:"max_minutes" binding:"required"`
}

func GenerateRouteHandler(c *gin.Context) {
	var req RouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Fetch points around the starting area
	targets, err := wikidata.FetchTargets(req.StartLat, req.StartLon, 3.0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed fetching wikidata"})
		return
	}

	// 2. Run the greedy path optimization
	routeResult := routing.GenerateRoute(req.StartLat, req.StartLon, req.EndLat, req.EndLon, req.MaxMinutes, targets)

	c.JSON(http.StatusOK, routeResult)
}

func GetTargetsHandler(c *gin.Context) {
	var params struct {
		Lat    float64 `form:"lat" binding:"required"`
		Lon    float64 `form:"lon" binding:"required"`
		Radius float64 `form:"radius"`
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'lat' or 'lon'"})
		return
	}

	if params.Radius == 0 {
		params.Radius = 2.0
	}

	targets, err := wikidata.FetchTargets(params.Lat, params.Lon, params.Radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":  len(targets),
		"target": targets,
	})
}

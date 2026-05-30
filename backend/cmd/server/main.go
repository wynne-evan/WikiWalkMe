package main

import (
	"log"
	"net/http"
	"time"
	"wikiwalkme-backend/internal/api"
	"wikiwalkme-backend/internal/routing"
	"wikiwalkme-backend/internal/wikidata"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No env file found, falling back to default")
	}

	wikiClient := wikidata.NewClient(30 * time.Minute)
	routeCache := routing.NewRouteCache(30 * time.Minute)
	apiCtx := &api.APIContext{
		WikiClient: wikiClient,
		RouteCache: routeCache,
	}

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	r.POST("/api/targets", apiCtx.GetTargetsHandler)
	r.POST("/api/route", apiCtx.GenerateRouteHandler)
	r.Run(":8080")
}

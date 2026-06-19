package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type checklistItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Done     bool   `json:"done"`
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	api.GET("/checklist", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"items": defaultChecklist()})
	})

	return router
}

func defaultChecklist() []checklistItem {
	return []checklistItem{
		{ID: "fuel", Title: "Fuel and fluids checked", Category: "Car"},
		{ID: "tires", Title: "Tire pressures set", Category: "Car"},
		{ID: "helmet", Title: "Helmet and gloves packed", Category: "Driver"},
		{ID: "license", Title: "License and registration ready", Category: "Admin"},
		{ID: "timing", Title: "Timing transponder charged", Category: "Gear"},
	}
}

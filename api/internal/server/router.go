package server

import (
	"context"
	"net/http"

	"raceday-checklist/api/internal/checklist"

	"github.com/gin-gonic/gin"
)

type ChecklistService interface {
	ListItems(ctx context.Context) ([]checklist.Item, error)
}

func NewRouter(checklistService ChecklistService) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	api.GET("/checklist", func(c *gin.Context) {
		items, err := checklistService.ListItems(c.Request.Context())
		if err != nil {
			if checklist.IsUnavailable(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "checklist unavailable"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	return router
}

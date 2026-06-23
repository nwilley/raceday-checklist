package server

import (
	"context"
	"net/http"

	"raceday-checklist/api/internal/checklist"

	"github.com/gin-gonic/gin"
)

const raceDayClientHeader = "X-Raceday-Client"

type ChecklistService interface {
	ListItems(ctx context.Context, clientID string) ([]checklist.Item, error)
	UpdateItemCompletion(ctx context.Context, clientID string, update checklist.CompletionUpdate) error
}

type completionRequest struct {
	Done bool `json:"done"`
}

func NewRouter(checklistService ChecklistService) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	api.GET("/checklist", func(c *gin.Context) {
		items, err := checklistService.ListItems(c.Request.Context(), c.GetHeader(raceDayClientHeader))
		if err != nil {
			writeChecklistError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	api.PATCH("/checklist/items/:sectionId/:itemId", func(c *gin.Context) {
		var request completionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid completion request"})
			return
		}

		err := checklistService.UpdateItemCompletion(c.Request.Context(), c.GetHeader(raceDayClientHeader), checklist.CompletionUpdate{
			SectionID: c.Param("sectionId"),
			ItemID:    c.Param("itemId"),
			Done:      request.Done,
		})
		if err != nil {
			writeChecklistError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"done": request.Done})
	})

	return router
}

func writeChecklistError(c *gin.Context, err error) {
	if checklist.IsInvalidClientID(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing race-day client id"})
		return
	}
	if checklist.IsUnavailable(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "checklist unavailable"})
		return
	}
	if checklist.IsNotFound(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "checklist item not found"})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

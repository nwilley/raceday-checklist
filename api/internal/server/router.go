package server

import (
	"context"
	"net/http"
	"os"
	"strings"

	"raceday-checklist/api/internal/checklist"

	"github.com/gin-gonic/gin"
)

const raceDayClientHeader = "X-Raceday-Client"
const corsAllowedOriginsEnv = "CORS_ALLOWED_ORIGINS"

type ChecklistService interface {
	ListItems(ctx context.Context, clientID string) ([]checklist.Item, error)
	UpdateItemCompletion(ctx context.Context, clientID string, update checklist.CompletionUpdate) error
}

type completionRequest struct {
	Done bool `json:"done"`
}

func NewRouter(checklistService ChecklistService) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware())

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

func corsMiddleware() gin.HandlerFunc {
	allowedOrigins := parseAllowedOrigins(os.Getenv(corsAllowedOriginsEnv))

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && isOriginAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, PATCH, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, "+raceDayClientHeader)
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func parseAllowedOrigins(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{"*"}
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	return false
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

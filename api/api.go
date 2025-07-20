package api

import (
	"github.com/gin-gonic/gin"
	"stravamcp/service"
)

func SetupRouter(activityService service.ActivityService) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	activityController := NewActivityController(activityService)
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/activities/refresh", activityController.RefreshActivities)
		apiGroup.GET("/activities", activityController.GetAllActivities)
		apiGroup.GET("/activities/:filter", activityController.GetAllActivities)
		apiGroup.GET("/activities/stream/:id", activityController.GetActivityStream)
	}
	return r
}

package api

import (
	"github.com/gin-gonic/gin"
	"stravamcp/service"
	"time"
)

type ActivityController interface {
	RefreshActivities(c *gin.Context)
	GetAllActivities(c *gin.Context)
	GetActivityStream(c *gin.Context)
}
type activityController struct {
	activityService service.ActivityService
}

func NewActivityController(activityService service.ActivityService) ActivityController {
	return &activityController{activityService: activityService}
}

func (ctrl *activityController) RefreshActivities(c *gin.Context) {
	err := ctrl.activityService.ProcessActivities(time.Now().Add(time.Hour * 24 * -365))
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.Status(200)
}

func (ctrl *activityController) GetAllActivities(c *gin.Context) {
	filter := c.Param("filter")
	activities, err := ctrl.activityService.GetAllActivities(c, filter, nil, nil)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, activities)
}

func (ctrl *activityController) GetActivityStream(c *gin.Context) {
	id := c.Param("id")
	activityStream, err := ctrl.activityService.GetActivityStream(c, id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, activityStream)
}

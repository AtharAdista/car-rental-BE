package handler

import (
	"carrental/internal/service/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DriverIncentiveV2Handler struct {
	driverIncentiveV2Service *service.DriverIncentiveV2Service
}

func NewDriverIncentiveV2Handler(driverIncentiveV2Service *service.DriverIncentiveV2Service) *DriverIncentiveV2Handler {
	return &DriverIncentiveV2Handler{driverIncentiveV2Service: driverIncentiveV2Service}
}

func (h *DriverIncentiveV2Handler) GetAllDriverIncentives(c *gin.Context) {

	driverIncentives, err := h.driverIncentiveV2Service.GetAllDriverIncentives()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": driverIncentives,
	})
}

func (h *DriverIncentiveV2Handler) GetDriverIncentiveById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	driverIncentive, err := h.driverIncentiveV2Service.GetDriverIncentiveById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": driverIncentive,
	})
}
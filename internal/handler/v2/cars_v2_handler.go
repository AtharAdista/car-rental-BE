package handler

import (
	"carrental/internal/model/v2"
	"carrental/internal/service/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CarsV2Handler struct {
	carsV2Service *service.CarsV2Service
}

func NewCarsV2Handler(carsV2Service *service.CarsV2Service) *CarsV2Handler {
	return &CarsV2Handler{carsV2Service: carsV2Service}
}

func (h *CarsV2Handler) CreateCar(c *gin.Context) {

	var req *model.CreateCarV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err := h.carsV2Service.CreateCar(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"name":       req.Name,
		"stock":      req.Stock,
		"daily_rent": req.DailyRent,
	})
}

func (h *CarsV2Handler) GetAllCars(c *gin.Context) {

	cars, err := h.carsV2Service.GetAllCars()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cars,
	})
}

func (h *CarsV2Handler) GetCarById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	car, err := h.carsV2Service.GetCarById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": car,
	})
}

func (h *CarsV2Handler) UpdateCarById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req *model.UpdateCarV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	car, err := h.carsV2Service.UpdateCarById(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": car})
}

func (h *CarsV2Handler) DeleteAllCars(c *gin.Context) {

	cars, err := h.carsV2Service.DeleteAllCars()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cars})
}

func (h *CarsV2Handler) DeleteCarById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	car, err := h.carsV2Service.DeleteCarById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": car})
}

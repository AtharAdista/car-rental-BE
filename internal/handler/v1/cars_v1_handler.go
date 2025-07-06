package handler

import (
	"carrental/internal/model/v1"
	"carrental/internal/service/v1"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CarsV1Handler struct {
	carsV1Service *service.CarsV1Service
}

func NewCarsV1Handler(carsV1Service *service.CarsV1Service) *CarsV1Handler {
	return &CarsV1Handler{carsV1Service: carsV1Service}
}

func (h *CarsV1Handler) CreateCar(c *gin.Context) {

	var req *model.CreateCarV1Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err := h.carsV1Service.CreateCar(req)

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

func (h *CarsV1Handler) GetAllCars(c *gin.Context) {

	cars, err := h.carsV1Service.GetAllCars()

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

func (h *CarsV1Handler) GetCarById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	car, err := h.carsV1Service.GetCarById(id)

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

func (h *CarsV1Handler) UpdateCarById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req *model.UpdateCarV1Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	car, err := h.carsV1Service.UpdateCarById(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": car})
}

func (h *CarsV1Handler) DeleteAllCars(c *gin.Context) {

	cars, err := h.carsV1Service.DeleteAllCars()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cars})
}

func (h *CarsV1Handler) DeleteCarById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	car, err := h.carsV1Service.DeleteCarById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": car})
}

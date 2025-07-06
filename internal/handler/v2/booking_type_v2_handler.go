package handler

import (
	"carrental/internal/model/v2"
	"carrental/internal/service/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingTypeV2Handler struct {
	bookingTypeService *service.BookingTypeV2Service
}

func NewBookingTypeV2Handler(bookingTypeV2Service *service.BookingTypeV2Service) *BookingTypeV2Handler {
	return &BookingTypeV2Handler{bookingTypeService: bookingTypeV2Service}
}

func (h *BookingTypeV2Handler) CreateBookingType(c *gin.Context) {
	var req *model.CreateBookingTypeV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id, err := h.bookingTypeService.CreateBookingType(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id": id,
		},
	})
}

func (h *BookingTypeV2Handler) GetAllBookingTypes(c *gin.Context) {

	bookingTypes, err := h.bookingTypeService.GetAllBookingTypes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookingTypes,
	})
}

func (h *BookingTypeV2Handler) GetBookingTypeById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	bookingType, err := h.bookingTypeService.GetBookingTypeById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookingType,
	})
}

func (h *BookingTypeV2Handler) UpdateBookingTypeById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req *model.UpdateBookingTypeV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	bookingType, err := h.bookingTypeService.UpdateCarById(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookingType})
}

func (h *BookingTypeV2Handler) DeleteAllBookingTypes(c *gin.Context) {

	bookingTypes, err := h.bookingTypeService.DeleteAllBookingTypes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookingTypes})
}

func (h *BookingTypeV2Handler) DeleteBookingTypeById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	bookingType, err := h.bookingTypeService.DeleteBookingTypeById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookingType})
}

package handler

import (
	"carrental/internal/model/v1"
	"carrental/internal/service/v1"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingV1Handler struct {
	bookingV1Service *service.BookingV1Service
}

func NewBookingV1Handler(bookingV1Service *service.BookingV1Service) *BookingV1Handler {
	return &BookingV1Handler{bookingV1Service: bookingV1Service}
}

func (h *BookingV1Handler) CreateBooking(c *gin.Context) {

	var req *model.CreateBookingV1Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id, err := h.bookingV1Service.CreateBooking(req)

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

func (h *BookingV1Handler) GetAllBookings(c *gin.Context) {

	booking, err := h.bookingV1Service.GetAllBookings()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": booking,
	})
}

func (h *BookingV1Handler) GetBookingById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	booking, err := h.bookingV1Service.GetBookingById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": booking,
	})
}

func (h *BookingV1Handler) UpdateBookingById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req *model.UpdateBookingV1Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	booking, err := h.bookingV1Service.UpdateBookingById(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func (h *BookingV1Handler) DeleteAllbookings(c *gin.Context) {

	bookings, err := h.bookingV1Service.DeleteAllBookings()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

func (h *BookingV1Handler) DeleteBookingById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	booking, err := h.bookingV1Service.DeleteBookingById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func (h *BookingV1Handler) FinishedStatusBooking(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	booking, err := h.bookingV1Service.FinishedStatusBooking(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

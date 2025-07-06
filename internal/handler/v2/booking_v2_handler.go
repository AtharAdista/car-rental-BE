package handler

import (
	"carrental/internal/model/v2"
	"carrental/internal/service/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingV2Handler struct {
	bookingV2Service *service.BookingV2Service
}

func NewBookingV2Handler(bookingV2Service *service.BookingV2Service) *BookingV2Handler {
	return &BookingV2Handler{bookingV2Service: bookingV2Service}
}

func (h *BookingV2Handler) CreateBooking(c *gin.Context) {

	var req *model.CreateBookingV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id, err := h.bookingV2Service.CreateBooking(req)

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

func (h *BookingV2Handler) GetAllBookings(c *gin.Context) {

	booking, err := h.bookingV2Service.GetAllBookings()

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

func (h *BookingV2Handler) GetBookingById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	booking, err := h.bookingV2Service.GetBookingById(id)

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

func (h *BookingV2Handler) UpdateBookingById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req *model.UpdateBookingV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	booking, err := h.bookingV2Service.UpdateBookingById(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func (h *BookingV2Handler) DeleteAllbookings(c *gin.Context) {

	bookings, err := h.bookingV2Service.DeleteAllBookings()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

func (h *BookingV2Handler) DeleteBookingById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	booking, err := h.bookingV2Service.DeleteBookingById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func (h *BookingV2Handler) FinishedStatusBooking(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	booking, err := h.bookingV2Service.FinishedStatusBooking(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

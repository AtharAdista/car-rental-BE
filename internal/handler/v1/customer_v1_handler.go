package handler

import (
	"carrental/internal/model/v1"
	"carrental/internal/service/v1"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerV1Handler struct {
	customerV1Service *service.CustomerV1Service
}

func NewCustomerV1Handler(customerV1Service *service.CustomerV1Service) *CustomerV1Handler {
	return &CustomerV1Handler{customerV1Service: customerV1Service}
}

func (h *CustomerV1Handler) CreateCustomer(c *gin.Context) {

	var req *model.CreateCustomerV1Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id, err := h.customerV1Service.CreateCustomer(req)

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

func (h *CustomerV1Handler) GetAllCustomers(c *gin.Context) {

	customers, err := h.customerV1Service.GetAllCustomers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": customers,
	})
}

func (h *CustomerV1Handler) GetCustomerById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	customer, err := h.customerV1Service.GetCustomerById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": customer,
	})
}

func (h *CustomerV1Handler) UpdateCustomerById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req *model.UpdateCustomerV1Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	customer, err := h.customerV1Service.UpdateCustomerById(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func (h *CustomerV1Handler) DeleteAllCustomers(c *gin.Context) {

	customers, err := h.customerV1Service.DeleteAllCustomers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

func (h *CustomerV1Handler) DeleteCustomerById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	customer, err := h.customerV1Service.DeleteCustomerById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

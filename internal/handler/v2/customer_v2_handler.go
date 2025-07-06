package handler

import (
	"carrental/internal/model/v2"
	"carrental/internal/service/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerV2Handler struct {
	customerV2Service *service.CustomerV2Service
}

func NewCustomerV2Handler(customerV2Service *service.CustomerV2Service) *CustomerV2Handler {
	return &CustomerV2Handler{customerV2Service: customerV2Service}
}

func (h *CustomerV2Handler) CreateCustomer(c *gin.Context) {

	var req *model.CreateCustomerV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id, err := h.customerV2Service.CreateCustomer(req)

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

func (h *CustomerV2Handler) GetAllCustomers(c *gin.Context) {

	customers, err := h.customerV2Service.GetAllCustomers()

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

func (h *CustomerV2Handler) GetCustomerById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	customer, err := h.customerV2Service.GetCustomerById(id)

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

func (h *CustomerV2Handler) UpdateCustomerById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req *model.UpdateCustomerV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	customer, err := h.customerV2Service.UpdateCustomerById(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func (h *CustomerV2Handler) DeleteAllCustomers(c *gin.Context) {

	customers, err := h.customerV2Service.DeleteAllCustomers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

func (h *CustomerV2Handler) DeleteCustomerById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	customer, err := h.customerV2Service.DeleteCustomerById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

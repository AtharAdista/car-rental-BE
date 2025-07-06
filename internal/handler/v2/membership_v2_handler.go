package handler

import (
	"carrental/internal/model/v2"
	"carrental/internal/service/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MembershipV2Handler struct {
	membershipV2Service *service.MembershipV2Service
}

func NewMembershipV2Handler(membershipV2Service *service.MembershipV2Service) *MembershipV2Handler {
	return &MembershipV2Handler{membershipV2Service: membershipV2Service}
}

func (h *MembershipV2Handler) CreateMembership(c *gin.Context) {
	var req *model.CreateMembershipV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id, err := h.membershipV2Service.CreateMembership(req)

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

func (h *MembershipV2Handler) GetAllMemberships(c *gin.Context) {

	memberships, err := h.membershipV2Service.GetAllMemberships()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": memberships,
	})
}

func (h *MembershipV2Handler) GetMembershipById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	membership, err := h.membershipV2Service.GetMembershipById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": membership,
	})
}

func (h *MembershipV2Handler) UpdateMembershipById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req *model.UpdateMembershipV2Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	membership, err := h.membershipV2Service.UpdateMembershipById(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": membership})
}

func (h *MembershipV2Handler) DeleteAllMemberships(c *gin.Context) {

	memberships, err := h.membershipV2Service.DeleteAllMembership()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": memberships})
}

func (h *MembershipV2Handler) DeleteMembershipById(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	Membership, err := h.membershipV2Service.DeleteMembershipById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": Membership})
}

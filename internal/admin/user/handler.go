package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"appsite-go/internal/apis/response"
	"appsite-go/internal/services/user/account"
	"appsite-go/internal/services/user/dto"
)

type Handler struct {
	svc *account.AuthService
}

func NewHandler(svc *account.AuthService) *Handler {
	return &Handler{svc: svc}
}

// ListUsers - Admin list users with full filters
func (h *Handler) ListUsers(c *gin.Context) {
	var req dto.UserFilterReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err)
		return
	}

	page := 1
	pageSize := 20
	
    if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	users, count, err := h.svc.ListUsers(req, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  users,
		"total": count,
	})
}

// GetUserDetail
func (h *Handler) GetUserDetail(c *gin.Context) {
	uid := c.Param("id")
	user, err := h.svc.GetDetail(uid)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, user)
}

// UpdateUser - generic update (ban, change group, etc.)
func (h *Handler) UpdateUser(c *gin.Context) {
	uid := c.Param("id")
	var req dto.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.svc.Update(uid, req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

package account

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"appsite-go/internal/apis/middleware"
	"appsite-go/internal/apis/response"
	apperr "appsite-go/internal/core/error"
	"appsite-go/internal/services/user/account"
	"appsite-go/internal/services/user/dto"
)

type Handler struct {
	svc *account.AuthService
}

func NewHandler(svc *account.AuthService) *Handler {
	return &Handler{svc: svc}
}

// GetProfile returns current user detail
func (h *Handler) GetProfile(c *gin.Context) {
	uid := c.GetString(middleware.ContextUserID)
	if uid == "" {
		response.Error(c, &apperr.AppError{Code: apperr.Unauthorized, Message: "unauthorized"})
		return
	}

	user, err := h.svc.GetDetail(uid)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// UpdateProfile updates user info
func (h *Handler) UpdateProfile(c *gin.Context) {
	uid := c.GetString(middleware.ContextUserID)
	if uid == "" {
		response.Error(c, &apperr.AppError{Code: apperr.Unauthorized, Message: "unauthorized"})
		return
	}

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

// ListUsers retrieves users based on filter
func (h *Handler) ListUsers(c *gin.Context) {
	var req dto.UserFilterReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Pagination
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

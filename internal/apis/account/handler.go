package account

import (
	"github.com/gin-gonic/gin"

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

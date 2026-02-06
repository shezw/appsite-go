package auth

import (
"github.com/gin-gonic/gin"

"appsite-go/internal/apis/response"
"appsite-go/internal/services/user/account"
)

// Handler handles admin authentication
type Handler struct {
svc *account.AuthService
}

// NewHandler creates a new admin auth handler
func NewHandler(svc *account.AuthService) *Handler {
return &Handler{svc: svc}
}

// LoginRequest represents admin login payload
type LoginRequest struct {
Username string `json:"username" binding:"required"`
Password string `json:"password" binding:"required"`
}

// Login handles admin login
func (h *Handler) Login(c *gin.Context) {
var req LoginRequest
if err := c.ShouldBindJSON(&req); err != nil {
response.Error(c, err)
return
}

token, user, err := h.svc.Login(req.Username, req.Password)
if err != nil {
response.Error(c, err)
return
}

// Strict check: Admin should enforce permissions.
// For now we assume if they can login here, we will check RBAC in middleware later.

response.Success(c, gin.H{
"token": token,
"admin": user,
})
}

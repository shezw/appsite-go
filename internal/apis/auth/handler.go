package auth

import (
"github.com/gin-gonic/gin"

"appsite-go/internal/apis/response"
"appsite-go/internal/services/user/account"
)

// Handler exposes auth operations
type Handler struct {
svc *account.AuthService
}

// NewHandler creates a new auth handler
func NewHandler(svc *account.AuthService) *Handler {
return &Handler{svc: svc}
}

// RegisterRequest represents the JSON body for registration
type RegisterRequest struct {
Username string `json:"username" binding:"required"`
Password string `json:"password" binding:"required,min=6"`
Email    string `json:"email" binding:"required,email"`
Mobile   string `json:"mobile"`
Nickname string `json:"nickname"`
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
var req RegisterRequest
if err := c.ShouldBindJSON(&req); err != nil {
// In a real app, you might reconstruct a validation error using apperr
response.Error(c, err)
return
}

user, err := h.svc.Register(account.RegisterInput{
Username: req.Username,
Password: req.Password,
Email:    req.Email,
Mobile:   req.Mobile,
Nickname: req.Nickname,
})
if err != nil {
response.Error(c, err)
return
}

response.Success(c, user)
}

// LoginRequest represents logic payload
type LoginRequest struct {
Identifier string `json:"identifier" binding:"required"`
Password   string `json:"password" binding:"required"`
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
var req LoginRequest
if err := c.ShouldBindJSON(&req); err != nil {
response.Error(c, err)
return
}

token, user, err := h.svc.Login(req.Identifier, req.Password)
if err != nil {
response.Error(c, err)
return
}

response.Success(c, gin.H{
"token": token,
"user":  user,
})
}

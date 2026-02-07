package apis

import (
	"github.com/gin-gonic/gin"

	"appsite-go/internal/apis/account"
	"appsite-go/internal/apis/auth"
	"appsite-go/internal/apis/middleware"
	"appsite-go/internal/services/access/token"
	account_svc "appsite-go/internal/services/user/account"
)

// Container holds all service dependencies for the API layer
type Container struct {
	TokenSvc *token.Service
	AuthSvc  *account_svc.AuthService
	// Add other services here as we implement their handlers
}

// RegisterRoutes registers all API routes
func RegisterRoutes(r *gin.Engine, c *Container) {
	v1 := r.Group("/api/v1")

	// Auth Routes (Public)
	if c.AuthSvc != nil {
		h := auth.NewHandler(c.AuthSvc)
		g := v1.Group("/auth")
		{
			g.POST("/register", h.Register)
			g.POST("/login", h.Login)
		}
	}
	
	// Account Routes (Protected)
	if c.AuthSvc != nil && c.TokenSvc != nil {
		h := account.NewHandler(c.AuthSvc)
		g := v1.Group("/account")
		g.Use(middleware.AuthMiddleware(c.TokenSvc))
		{
			g.GET("/profile", h.GetProfile)
			g.PUT("/profile", h.UpdateProfile)
		}
	}
}

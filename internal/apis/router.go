package apis

import (
	"github.com/gin-gonic/gin"

	"appsite-go/internal/apis/auth"
	"appsite-go/internal/services/user/account"
)

// Container holds all service dependencies for the API layer
type Container struct {
	AuthSvc *account.AuthService
	// Add other services here as we implement their handlers
}

// RegisterRoutes registers all API routes
func RegisterRoutes(r *gin.Engine, c *Container) {
	v1 := r.Group("/api/v1")

	// Auth Routes
	if c.AuthSvc != nil {
		h := auth.NewHandler(c.AuthSvc)
		g := v1.Group("/auth")
		{
			g.POST("/register", h.Register)
			g.POST("/login", h.Login)
		}
	}
}

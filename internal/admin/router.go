package admin

import (
	"github.com/gin-gonic/gin"

	"appsite-go/internal/admin/auth"
	"appsite-go/internal/admin/user"
	"appsite-go/internal/services/user/account"
)

// Container holds dependencies for admin handlers
type Container struct {
	AuthSvc *account.AuthService
}

// RegisterRoutes registers admin routes
func RegisterRoutes(r *gin.Engine, c *Container) {
	v1 := r.Group("/admin/v1")

	// Auth
	if c.AuthSvc != nil {
		h := auth.NewHandler(c.AuthSvc)
		v1.POST("/login", h.Login)
	}

	// Users
	if c.AuthSvc != nil {
		h := user.NewHandler(c.AuthSvc)
		// TODO: Add Admin Middleware here
		g := v1.Group("/users")
		{
			g.GET("", h.ListUsers)
			g.GET("/:id", h.GetUserDetail)
			g.PUT("/:id", h.UpdateUser)
		}
	}
}

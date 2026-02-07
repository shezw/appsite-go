package admin

import (
	"github.com/gin-gonic/gin"

	"appsite-go/internal/admin/auth"
	"appsite-go/internal/admin/contents"
	"appsite-go/internal/admin/system"
	"appsite-go/internal/admin/user"
	"appsite-go/internal/services/user/account"
	scontent "appsite-go/internal/services/contents"
	"appsite-go/internal/core/setting"
)

// Container holds dependencies for admin handlers
type Container struct {
	AuthSvc    *account.AuthService
	ArticleSvc *scontent.ArticleService
	BannerSvc  *scontent.BannerService
	Config     *setting.Config
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

	// Content
	if c.ArticleSvc != nil && c.BannerSvc != nil {
		h := contents.NewHandler(c.ArticleSvc, c.BannerSvc)
		g := v1.Group("/contents") 
		{
			// Articles
			g.GET("/articles", h.ListArticles)
			g.POST("/articles", h.CreateArticle)
			g.GET("/articles/:id", h.GetArticle)
			g.PUT("/articles/:id", h.UpdateArticle)
			g.DELETE("/articles/:id", h.DeleteArticle)
			
			// Banners
			g.GET("/banners", h.ListBanners)
			g.POST("/banners", h.CreateBanner)
			g.GET("/banners/:id", h.GetBanner)
			g.PUT("/banners/:id", h.UpdateBanner)
			g.DELETE("/banners/:id", h.DeleteBanner)
		}
	}

	// System & Config
	if c.Config != nil {
		h := system.NewHandler(c.Config)
		v1.GET("/menu", h.GetMenu)
	}
}

package apis

import (
	"github.com/gin-gonic/gin"

	"appsite-go/internal/apis/account"
	"appsite-go/internal/apis/auth"
	"appsite-go/internal/apis/content"
	"appsite-go/internal/apis/middleware"
	"appsite-go/internal/services/access/token"
	"appsite-go/internal/services/contents"
	account_svc "appsite-go/internal/services/user/account"
)

// Container holds all service dependencies for the API layer
type Container struct {
	TokenSvc   *token.Service
	AuthSvc    *account_svc.AuthService
	ArticleSvc *contents.ArticleService
	BannerSvc  *contents.BannerService
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

		// Users (Admin/Public Directory)
		u := v1.Group("/users")
		u.Use(middleware.AuthMiddleware(c.TokenSvc))
		{
			u.GET("", h.ListUsers)
		}
	}

	// Content Routes (Public Read, Protected Write)
	if c.ArticleSvc != nil && c.BannerSvc != nil {
		h := content.NewHandler(c.ArticleSvc, c.BannerSvc)
		
		g := v1.Group("/content")
		{
			g.GET("/articles", h.ListArticles)
			g.GET("/articles/:id", h.GetArticle)
			g.GET("/banners", h.ListBanners)
		}
		
		// Protected
		p := v1.Group("/content")
		if c.TokenSvc != nil {
			p.Use(middleware.AuthMiddleware(c.TokenSvc))
		}
		{
			p.POST("/articles", h.CreateArticle)
			p.PUT("/articles/:id", h.UpdateArticle)
			p.DELETE("/articles/:id", h.DeleteArticle)
			p.POST("/banners", h.CreateBanner)
		}
	}
}

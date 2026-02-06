package route

import (
	"github.com/gin-gonic/gin"
	"appsite-go/internal/core/setting"
)

// NewEngine initializes the Gin engine with global middlewares
func NewEngine(cfg *setting.Config) *gin.Engine {
	// Set Gin mode
	switch cfg.App.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Apply Core Middlewares
	r.Use(LoggerMiddleware())
	r.Use(RecoveryMiddleware())
	r.Use(CORSMiddleware())
	r.Use(SaasMiddleware())

	return r
}

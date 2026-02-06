package route_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"appsite-go/internal/core/route"
	"appsite-go/internal/core/setting"
)

func TestNewEngine(t *testing.T) {
	cfg := &setting.Config{App: setting.AppConfig{Mode: "test"}}
	r := route.NewEngine(cfg)

	// Add a ping route
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestSaasMiddleware(t *testing.T) {
	// Middleware is internal/core/route, we can test it indirectly via Engine or unit test it if exported
	// SaasMiddleware is exported.
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(route.SaasMiddleware())
	
	r.GET("/tenant", func(c *gin.Context) {
		tid, _ := c.Get("tenant_id")
		c.String(200, tid.(string))
	})
	
	// 1. Header
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tenant", nil)
	req.Header.Set("X-Tenant-ID", "tenant-123")
	r.ServeHTTP(w, req)
	
	if w.Body.String() != "tenant-123" {
		t.Errorf("Expected tenant-123, got %s", w.Body.String())
	}
}

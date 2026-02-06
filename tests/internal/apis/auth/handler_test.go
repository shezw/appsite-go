package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/apis/auth"
	"appsite-go/internal/core/setting"
	"appsite-go/internal/services/access/token"
	"appsite-go/internal/services/access/verify"
	"appsite-go/internal/services/user/account"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func setupAuthService(t *testing.T) *account.AuthService {
	db := setupDB(t)
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)

	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	otpSvc := verify.NewOTPService(rdb)

	cfg := setting.AppConfig{
		JwtSecret: "test_secret",
		JwtExpire: time.Hour,
	}
	tokenSvc := token.NewService(cfg)
	return account.NewAuthService(db, tokenSvc, otpSvc)
}

func TestRegister(t *testing.T) {
	svc := setupAuthService(t)
	handler := auth.NewHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", handler.Register)

	body := map[string]string{
		"username": "api_user",
		"password": "password123",
		"email":    "api@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	// You can parse Body and check "code": 200
}

func TestLogin(t *testing.T) {
	svc := setupAuthService(t)
	handler := auth.NewHandler(svc)
	
	// Helper to create user
	_, _ = svc.Register(account.RegisterInput{
		Username: "login_user",
		Password: "password123",
		Email:    "login@example.com",
	})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", handler.Login)

	body := map[string]string{
		"identifier": "login_user",
		"password":   "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

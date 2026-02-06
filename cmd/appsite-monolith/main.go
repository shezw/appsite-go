package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"appsite-go/internal/apis"
	"appsite-go/internal/core/log"
	"appsite-go/internal/core/route"
	"appsite-go/internal/core/setting"
	"appsite-go/internal/services/access/token"
	"appsite-go/internal/services/access/verify"
	"appsite-go/internal/services/user/account"
	"appsite-go/pkg/utils/orm"
	"appsite-go/pkg/utils/redis"
)

func main() {
	// 1. Load Config
	loader, err := setting.NewLoader("configs", "config", "yaml")
	if err != nil {
		fmt.Printf("Failed to init config loader: %v\n", err)
		os.Exit(1)
	}
	cfg, err := loader.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Logger
	zapLogger, err := log.NewZapLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	log.SetLogger(zapLogger)
	defer zapLogger.Sync()

	log.Info(context.Background(), "Starting Appsite Monolith...")

	// 3. Initialize Database
	db, err := orm.NewMySQLConnection(&cfg.Database)
	if err != nil {
		log.Fatal(context.Background(), "Failed to connect to database", "err", err)
	}
	log.Info(context.Background(), "Database connected")

	// 4. Initialize Redis
	rdb, err := redis.NewClient(&cfg.Redis)
	if err != nil {
		log.Fatal(context.Background(), "Failed to connect to redis", "err", err)
	}
	log.Info(context.Background(), "Redis connected")

	// 5. Initialize Services (DI)
	// Access Services
	tokenSvc := token.NewService(cfg.App)
	otpSvc := verify.NewOTPService(rdb)

	// User Services
	authSvc := account.NewAuthService(db, tokenSvc, otpSvc)

	// ... Init other services here ...

	// 6. Initialize API Container
	container := &apis.Container{
		AuthSvc: authSvc,
		// ... Inject other services ...
	}

	// 7. Setup Router
	r := route.NewEngine(cfg)
	apis.RegisterRoutes(r, container)
	// admin.RegisterRoutes(r, adminContainer) 

	// 8. Run Server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		log.Info(context.Background(), "Server listening", "addr", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(context.Background(), "Listen error", "err", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info(context.Background(), "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(context.Background(), "Server forced to shutdown", "err", err)
	}

	log.Info(context.Background(), "Server exiting")
}

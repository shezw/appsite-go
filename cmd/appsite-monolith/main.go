package main

import (
"context"
"fmt"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"github.com/alicebob/miniredis/v2"
goredis "github.com/redis/go-redis/v9"

"appsite-go/internal/apis"
"appsite-go/internal/core/log"
"appsite-go/internal/core/route"
"appsite-go/internal/core/setting"
"appsite-go/internal/services/access/token"
"appsite-go/internal/services/access/verify"
"appsite-go/internal/services/user/account"
"appsite-go/pkg/utils/orm"
appsite_redis "appsite-go/pkg/utils/redis"
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

ctx := context.Background()
log.Info(ctx, "Starting Appsite Monolith...")

// 3. Initialize Database
db, err := orm.InitDB(&cfg.Database)
if err != nil {
log.Fatal(ctx, "Failed to connect to database", "err", err)
}
log.Info(ctx, "Database connected")

// 4. Initialize Redis
// Try connecting to configured Redis
var rdb *goredis.Client
rdb, err = appsite_redis.NewClient(&cfg.Redis)
if err != nil {
log.Warn(ctx, "Failed to connect to configured Redis. Falling back to embedded Miniredis.", "err", err)

// Fallback to Miniredis
mr, errMr := miniredis.Run()
if errMr != nil {
log.Fatal(ctx, "Failed to start embedded Miniredis", "err", errMr)
}
// Miniredis address
rdb = goredis.NewClient(&goredis.Options{
Addr: mr.Addr(),
})
log.Info(ctx, "Miniredis started", "addr", mr.Addr())
} else {
log.Info(ctx, "Redis connected")
}

// 5. Initialize Services (DI)
// Access Services
tokenSvc := token.NewService(cfg.App)
otpSvc := verify.NewOTPService(rdb)

// User Services
authSvc := account.NewAuthService(db, tokenSvc, otpSvc)

// ... Init other services here ...

// 6. Initialize API Container
	// 6. Initialize API Container
	container := &apis.Container{
		TokenSvc: tokenSvc,
		AuthSvc:  authSvc,
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
log.Info(ctx, "Server listening", "addr", serverAddr)
if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
log.Fatal(ctx, "Listen error", "err", err)
}
}()

// Graceful Shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
log.Info(ctx, "Shutting down server...")

shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err := srv.Shutdown(shutdownCtx); err != nil {
log.Fatal(ctx, "Server forced to shutdown", "err", err)
}

log.Info(ctx, "Server exiting")
}

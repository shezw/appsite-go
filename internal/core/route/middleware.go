package route

import (
"time"

"github.com/gin-gonic/gin"


"appsite-go/internal/core/log"
)

const (
HeaderTenantID = "X-Tenant-ID"
ContextTenantID = "tenant_id"
)

// LoggerMiddleware logs valid inquiries
func LoggerMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
start := time.Now()
path := c.Request.URL.Path
query := c.Request.URL.RawQuery

c.Next()

cost := time.Since(start)
status := c.Writer.Status()

// Prepare fields as Key-Value pairs for the logger abstraction
fields := []interface{}{
"status", status,
"method", c.Request.Method,
"path", path,
"query", query,
"ip", c.ClientIP(),
"user-agent", c.Request.UserAgent(),
"cost", cost,
}

ctx := c.Request.Context()

if len(c.Errors) > 0 {
for _, e := range c.Errors.Errors() {
// Error message plus fields
log.Error(ctx, e, fields...)
}
} else {
if status >= 400 {
log.Warn(ctx, "HTTP Request", fields...)
} else {
log.Info(ctx, "HTTP Request", fields...)
}
}
}
}

// RecoveryMiddleware recovers from panic
func RecoveryMiddleware() gin.HandlerFunc {
return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
ctx := c.Request.Context()
if err, ok := recovered.(error); ok {
log.Error(ctx, "Panic Recovered", "error", err)
} else {
log.Error(ctx, "Panic Recovered", "error", recovered)
}
c.AbortWithStatus(500)
})
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-ID")
c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

if c.Request.Method == "OPTIONS" {
c.AbortWithStatus(204)
return
}

c.Next()
}
}

// SaasMiddleware extracts tenant information
func SaasMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
tenantID := c.GetHeader(HeaderTenantID)
if tenantID == "" {
// Fallback: Check Query param
tenantID = c.Query("tenant_id")
}

if tenantID != "" {
c.Set(ContextTenantID, tenantID)
}

c.Next()
}
}

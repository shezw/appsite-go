package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"appsite-go/internal/apis/response"
	apperr "appsite-go/internal/core/error"
	"appsite-go/internal/services/access/token"
)

const (
	ContextUserID = "user_id"
	ContextUser   = "user_claims"
)

// AuthMiddleware verifies JWT token
func AuthMiddleware(svc *token.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, &apperr.AppError{Code: apperr.Unauthorized, Message: "unauthorized"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, &apperr.AppError{Code: apperr.Unauthorized, Message: "invalid token format"})
			c.Abort()
			return
		}

		claims, err := svc.ParseToken(parts[1])
		if err != nil {
			response.Error(c, &apperr.AppError{Code: apperr.Unauthorized, Message: "invalid token", Err: err})
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUser, claims)
		c.Next()
	}
}

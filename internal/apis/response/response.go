package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperr "appsite-go/internal/core/error"
)

// Response standardizes the JSON response format
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success sends a successful JSON response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: int(apperr.Success),
		Msg:  "success",
		Data: data,
	})
}

// Error sends an error JSON response
func Error(c *gin.Context, err error) {
	var e *apperr.AppError
	if errors.As(err, &e) {
		// We use 200 OK for business errors to allow frontend to parse the JSON body
		// unless it's a very specific protocol requirement.
		// You can also map specific apperr.Code to http status if preferred.
		c.JSON(http.StatusOK, Response{
			Code: int(e.Code),
			Msg:  e.Message,
		})
		return
	}

	// Unknown error
	c.JSON(http.StatusInternalServerError, Response{
		Code: int(apperr.ServerError),
		Msg:  err.Error(),
	})
}

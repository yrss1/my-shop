package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Object struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

func OK(c *gin.Context, data any) {
	h := Object{
		Success: true,
		Data:    data,
	}
	c.JSON(http.StatusOK, h)
}

func Created(c *gin.Context, data any) {
	h := Object{
		Success: true,
		Data:    data,
	}
	c.JSON(http.StatusCreated, h)
}

func BadRequest(c *gin.Context, err error, data any) {
	h := Object{
		Success: false,
		Message: err.Error(),
		Data:    data,
	}
	c.JSON(http.StatusBadRequest, h)
}

func NotFound(c *gin.Context, err error) {
	h := Object{
		Success: false,
		Message: err.Error(),
	}
	c.JSON(http.StatusNotFound, h)
}

func InternalServerError(c *gin.Context, err error) {
	h := Object{
		Success: false,
		Message: err.Error(),
	}
	c.JSON(http.StatusInternalServerError, h)
}

func MethodNotAllowedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedMethods := map[string]bool{
			http.MethodGet:    true,
			http.MethodPost:   true,
			http.MethodPut:    true,
			http.MethodDelete: true,
		}
		if !allowedMethods[c.Request.Method] {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method Not Allowed"})
			c.Abort()
			return
		}

		c.Next()
	}
}

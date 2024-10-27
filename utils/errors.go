package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HandleError(c *gin.Context, status int, message string) {
	Error(message, Fields(map[string]interface{}{
		"status": status,
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	}))

	c.JSON(status, ErrorResponse{
		Status:  status,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	HandleError(c, http.StatusBadRequest, message)
}

func InternalError(c *gin.Context, message string) {
	HandleError(c, http.StatusInternalServerError, message)
}

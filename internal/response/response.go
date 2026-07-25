package response

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	RequestID string   `json:"requestId"`
	Error     APIError `json:"error"`
}

func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, errorCode string, message string) {
	reqID := c.GetString("RequestID")
	c.JSON(statusCode, ErrorResponse{
		RequestID: reqID,
		Error: APIError{
			Code:    errorCode,
			Message: message,
		},
	})
}

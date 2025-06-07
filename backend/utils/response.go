package utils

import (
	"time"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
}

func ErrorResponse(message string, error interface{}) APIResponse {
	return APIResponse{
		Success:   false,
		Message:   message,
		Error:     error,
		Timestamp: time.Now().UTC(),
	}
}

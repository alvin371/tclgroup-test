package httpx

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

const apiVersion = "1.0.0"

type successEnvelope struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Success    bool        `json:"success"`
	Version    string      `json:"version"`
}

type listEnvelope struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Items      interface{} `json:"items"`
	Meta       metaInfo    `json:"meta"`
	Success    bool        `json:"success"`
	Version    string      `json:"version"`
}

type metaInfo struct {
	TotalPage int   `json:"total_page"`
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
}

type errorEnvelope struct {
	Success bool      `json:"success"`
	Error   *apiError `json:"error,omitempty"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success writes a 200 JSON response with data.
func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, successEnvelope{
		StatusCode: http.StatusOK,
		Message:    msg,
		Data:       data,
		Success:    true,
		Version:    apiVersion,
	})
}

// Created writes a 201 JSON response with data.
func Created(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusCreated, successEnvelope{
		StatusCode: http.StatusCreated,
		Message:    msg,
		Data:       data,
		Success:    true,
		Version:    apiVersion,
	})
}

// Error writes an error JSON response with the given HTTP status code.
func Error(c *gin.Context, status int, code, msg string) {
	c.JSON(status, errorEnvelope{
		Success: false,
		Error:   &apiError{Code: code, Message: msg},
	})
}

// Paginated writes a 200 JSON response with items and pagination metadata.
func Paginated(c *gin.Context, items interface{}, total int64, page, perPage int, msg string) {
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, listEnvelope{
		StatusCode: http.StatusOK,
		Message:    msg,
		Items:      items,
		Meta: metaInfo{
			TotalPage: totalPage,
			Total:     total,
			Page:      page,
			PerPage:   perPage,
		},
		Success: true,
		Version: apiVersion,
	})
}

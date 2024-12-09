package helper

import (
	"github.com/gin-gonic/gin"
)

type HTTPResponse struct {
	ErrorMsg   string      `json:"error_msg,omitempty"`
	Message    string      `json:"message,omitempty"`
	Page       int         `json:"page,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	TotalItems int         `json:"total_items,omitempty"`
	TotalPages int         `json:"total_pages,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func ResponseOK(c *gin.Context, data interface{}, message string, httpStatusCode int) {
	c.JSON(httpStatusCode, HTTPResponse{
		Message: message,
		Data:    data,
	})
}
func ResponseOKPagination(c *gin.Context, data interface{}, message string, page, limit, totalItems, totalPages, httpStatusCode int) {
	c.JSON(httpStatusCode, HTTPResponse{
		Message: message,
		Page: page,
		Limit: limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Data:    data,
	})
}

func ResponseError(c *gin.Context, errorMsg string, message string, httpStatusCode int) {
	c.JSON(httpStatusCode, HTTPResponse{
		ErrorMsg: errorMsg,
		Message:  message,
	})
}

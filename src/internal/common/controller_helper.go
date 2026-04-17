package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindJSON là helper để bind JSON và trả về lỗi 400 nếu thất bại
func BindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		IsErrorCustom(c, err, http.StatusBadRequest, "Invalid request body")
		return err
	}
	return nil
}

// IsErrorCustom trả về lỗi với status code và message tùy chỉnh
func IsErrorCustom(c *gin.Context, err error, status int, message string) {
	c.JSON(status, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}

// IsErrorWithMessage trả về lỗi 500 với message tùy chỉnh
func IsErrorWithMessage(c *gin.Context, err error, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}

// IsValidationError trả về true nếu là lỗi validation (ví dụ đơn giản)
func IsValidationError(err error) bool {
	// Trong thực tế bạn có thể check type của error (ví dụ validator.ValidationErrors)
	return err != nil
}

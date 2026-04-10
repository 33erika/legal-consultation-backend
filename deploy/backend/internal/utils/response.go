package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// 错误响应
func Error(c *gin.Context, httpStatus int, code string, message string) {
	c.JSON(httpStatus, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "ERR_BAD_REQUEST", message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "ERR_UNAUTHORIZED", message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "ERR_FORBIDDEN", message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "ERR_NOT_FOUND", message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "ERR_INTERNAL", message)
}

// 错误码常量
const (
	ErrAuthInvalidCredentials = "ERR_AUTH_INVALID_CREDENTIALS"
	ErrAuthTokenExpired       = "ERR_AUTH_TOKEN_EXPIRED"
	ErrAuthForbidden          = "ERR_AUTH_FORBIDDEN"

	ErrConsultationNotFound          = "ERR_CONSULTATION_NOT_FOUND"
	ErrConsultationAlreadyAccepted  = "ERR_CONSULTATION_ALREADY_ACCEPTED"
	ErrConsultationAlreadyClosed    = "ERR_CONSULTATION_ALREADY_CLOSED"
	ErrConsultationCannotTransfer   = "ERR_CONSULTATION_CANNOT_TRANSFER"

	ErrTemplateNotFound             = "ERR_TEMPLATE_NOT_FOUND"
	ErrTemplateAlreadyPublished     = "ERR_TEMPLATE_ALREADY_PUBLISHED"
	ErrTemplateRequestNotFound      = "ERR_TEMPLATE_REQUEST_NOT_FOUND"

	ErrUserNotFound       = "ERR_USER_NOT_FOUND"
	ErrDepartmentNotEmpty = "ERR_DEPARTMENT_NOT_EMPTY"

	ErrFileTypeNotAllowed = "ERR_FILE_TYPE_NOT_ALLOWED"
	ErrFileSizeExceeded   = "ERR_FILE_SIZE_EXCEEDED"

	ErrNotificationFailed = "ERR_NOTIFICATION_FAILED"
)

package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrCodeBadeRequest ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound    ErrorCode = "NOT_FOUND"
	ErrCodeConflict    ErrorCode = "CONFLICT"
	ErrCodeInternal    ErrorCode = "INTERNAL_SERVER_ERR"
)

type AppError struct {
	Message string
	Code    ErrorCode
	Err     error
}

type APIResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Pagination any    `json:"pagination,omitempty`
}

func (ae *AppError) Error() string {
	return ""
}

func NewError(message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

func WrapError(message string, code ErrorCode, err error) error {
	return &AppError{
		Message: message,
		Code:    code,
		Err:     err,
	}
}

func ResponseError(ctx *gin.Context, err error) {
	if appError, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appError.Code)
		response := gin.H{
			"error": CapitalizeFirst(appError.Message),
			"code":  appError.Code,
		}

		if appError.Err != nil {
			response["detail"] = appError.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}

func ResponseSuccess(ctx *gin.Context, status int, message string, data ...any) {
	resp := APIResponse{
		Status:  "success",
		Message: CapitalizeFirst(message),
	}

	if len(data) > 0 && data[0] != nil {
		resp.Data = data[0]
	}

	ctx.JSON(status, resp)
}

func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}

func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
}

func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadeRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

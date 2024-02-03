package pnd

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type AppErrorCode string

const (
	// Common errors
	ErrCodeBadRequest        AppErrorCode = "ERR_BAD_REQUEST"
	ErrCodeInvalidParam      AppErrorCode = "ERR_INVALID_PARAM"
	ErrCodeInvalidPagination AppErrorCode = "ERR_INVALID_PAGINATION"
	ErrCodeInvalidQuery      AppErrorCode = "ERR_INVALID_QUERY"
	ErrCodeInvalidBody       AppErrorCode = "ERR_INVALID_BODY"
	ErrCodeMultipartForm     AppErrorCode = "ERR_MULTIPART_FORM"

	// Common errors - Auth
	ErrCodeInvalidFBToken    AppErrorCode = "ERR_INVALID_FB_TOKEN"
	ErrCodeUserNotRegistered AppErrorCode = "ERR_USER_NOT_REGISTERED"
	ErrCodeForbidden         AppErrorCode = "ERR_FORBIDDEN"

	// Common Errors - Resource
	ErrCodeNotFound AppErrorCode = "ERR_NOT_FOUND"
	ErrCodeConflict AppErrorCode = "ERR_CONFLICT"

	ErrCodeUnknown AppErrorCode = "ERR_UNKNOWN"
)

type AppError struct {
	Err        error `json:"-"`
	StatusCode int   `json:"-"`

	Code    AppErrorCode `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
}

func (e *AppError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func (e *AppError) Error() string {
	return e.Message
}

func ErrCustom(err error, statusCode int, code AppErrorCode, message string) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func ErrBadRequest(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		Code:       ErrCodeBadRequest,
		Message:    err.Error(),
	}
}

func ErrInvalidParam(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		Code:       ErrCodeInvalidParam,
		Message:    err.Error(),
	}
}

func ErrInvalidPagination(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		Code:       ErrCodeInvalidPagination,
		Message:    err.Error(),
	}
}

func ErrInvalidQuery(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		Code:       ErrCodeInvalidQuery,
		Message:    err.Error(),
	}
}

func ErrInvalidBody(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		Code:       ErrCodeInvalidBody,
		Message:    err.Error(),
	}
}

func ErrMultipartFormError(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		Code:       ErrCodeMultipartForm,
		Message:    err.Error(),
	}
}

func ErrInvalidFBToken(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusUnauthorized,
		Code:       ErrCodeInvalidFBToken,
		Message:    err.Error(),
	}
}

func ErrUserNotRegistered(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusUnauthorized,
		Code:       ErrCodeUserNotRegistered,
		Message:    err.Error(),
	}
}

func ErrForbidden(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusForbidden,
		Code:       ErrCodeForbidden,
		Message:    err.Error(),
	}
}

func ErrNotFound(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusNotFound,
		Code:       ErrCodeNotFound,
		Message:    err.Error(),
	}
}

func ErrConflict(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusConflict,
		Code:       ErrCodeConflict,
		Message:    err.Error(),
	}
}

func ErrUnknown(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		Code:       ErrCodeUnknown,
		Message:    err.Error(),
	}
}

func FromPostgresError(err error) *AppError {
	errStr := err.Error()
	if strings.Contains(errStr, "no rows in result set") {
		return ErrCustom(err, http.StatusNotFound, ErrCodeNotFound, "해당하는 자원이 없습니다")
	} else if strings.Contains(errStr, "violates foreign key constraint") {
		return ErrCustom(err, http.StatusNotFound, ErrCodeNotFound, "해당하는 자원이 없습니다")
	} else if strings.Contains(errStr, "violates not-null constraint") {
		return ErrCustom(err, http.StatusBadRequest, ErrCodeBadRequest, "필수 값이 누락되었습니다")
	} else if strings.Contains(errStr, "violates check constraint") {
		return ErrCustom(err, http.StatusBadRequest, ErrCodeBadRequest, "잘못된 값입니다")
	} else if strings.Contains(errStr, "violates unique constraint") {
		return ErrCustom(err, http.StatusConflict, ErrCodeConflict, "중복된 값입니다")
	} else {
		return ErrCustom(err, http.StatusInternalServerError, ErrCodeUnknown, "알 수 없는 오류가 발생했습니다")
	}
}

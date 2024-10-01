package pnd

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
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
	ErrCodeInvalidFBToken     AppErrorCode = "ERR_INVALID_FB_TOKEN"     //nolint:gosec
	ErrCodeInvalidBearerToken AppErrorCode = "ERR_INVALID_BEARER_TOKEN" //nolint:gosec
	ErrCodeUserNotRegistered  AppErrorCode = "ERR_USER_NOT_REGISTERED"
	ErrCodeForbidden          AppErrorCode = "ERR_FORBIDDEN"

	// Common Errors - Resource
	ErrCodeNotFound AppErrorCode = "ERR_NOT_FOUND"
	ErrCodeConflict AppErrorCode = "ERR_CONFLICT"

	// Common Errors - Chat
	ErrCodeClientRegistrationFailed   AppErrorCode = "ERR_CLIENT_REGISTRATION_FAILED"
	ErrCodeClientUnregistrationFailed AppErrorCode = "ERR_CLIENT_UNREGISTRATION_FAILED"
	ErrCodeMessageEncodingFailed      AppErrorCode = "ERR_MESSAGE_ENCODING_FAILED"
	ErrCodeRoomCreationFailed         AppErrorCode = "ERR_ROOM_CREATION_FAILED"
	ErrCodeRoomNotFound               AppErrorCode = "ERR_ROOM_NOT_FOUND"

	ErrCodeUnknown AppErrorCode = "ERR_UNKNOWN"
)

type AppError struct {
	Err        error `json:"-"`
	StatusCode int   `json:"-"`

	Code    AppErrorCode `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
}

func NewAppError(err error, statusCode int, code AppErrorCode, message string) *AppError {
	log.Error().Err(err).Msg(message)
	return &AppError{
		Err:        err,
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func ErrDefault(err error, statusCode int, code AppErrorCode) *AppError {
	return NewAppError(err, statusCode, code, err.Error())
}

func ErrBadRequest(err error) *AppError {
	return ErrDefault(err, http.StatusBadRequest, ErrCodeBadRequest)
}

func ErrInvalidParam(err error) *AppError {
	return ErrDefault(err, http.StatusBadRequest, ErrCodeInvalidParam)
}

func ErrInvalidPagination(err error) *AppError {
	return ErrDefault(err, http.StatusBadRequest, ErrCodeInvalidPagination)
}

func ErrInvalidQuery(err error) *AppError {
	return ErrDefault(err, http.StatusBadRequest, ErrCodeInvalidQuery)
}

func ErrInvalidBody(err error) *AppError {
	return ErrDefault(err, http.StatusBadRequest, ErrCodeInvalidBody)
}

func ErrMultipartFormError(err error) *AppError {
	return ErrDefault(err, http.StatusBadRequest, ErrCodeMultipartForm)
}

func ErrInvalidFBToken(err error) *AppError {
	return ErrDefault(err, http.StatusUnauthorized, ErrCodeInvalidFBToken)
}

func ErrInvalidBearerToken(err error) *AppError {
	return ErrDefault(err, http.StatusUnauthorized, ErrCodeInvalidBearerToken)
}

func ErrUserNotRegistered(err error) *AppError {
	return ErrDefault(err, http.StatusUnauthorized, ErrCodeUserNotRegistered)
}

func ErrForbidden(err error) *AppError {
	return ErrDefault(err, http.StatusForbidden, ErrCodeForbidden)
}

func ErrNotFound(err error) *AppError {
	return ErrDefault(err, http.StatusNotFound, ErrCodeNotFound)
}

func ErrConflict(err error) *AppError {
	return ErrDefault(err, http.StatusConflict, ErrCodeConflict)
}

func ErrUnknown(err error) *AppError {
	return ErrDefault(err, http.StatusInternalServerError, ErrCodeUnknown)
}

func FromPostgresError(err error) *AppError {
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "no rows in result set"):
		return NewAppError(err, http.StatusNotFound, ErrCodeNotFound, "해당하는 자원이 없습니다")
	case strings.Contains(errStr, "violates foreign key constraint"):
		return NewAppError(err, http.StatusNotFound, ErrCodeNotFound, "해당하는 자원이 없습니다")
	case strings.Contains(errStr, "violates not-null constraint"):
		return NewAppError(err, http.StatusBadRequest, ErrCodeBadRequest, "필수 값이 누락되었습니다")
	case strings.Contains(errStr, "violates check constraint"):
		return NewAppError(err, http.StatusBadRequest, ErrCodeBadRequest, "잘못된 값입니다")
	case strings.Contains(errStr, "violates unique constraint"):
		return NewAppError(err, http.StatusConflict, ErrCodeConflict, "중복된 값입니다")
	default:
		return NewAppError(err, http.StatusInternalServerError, ErrCodeUnknown, "알 수 없는 오류가 발생했습니다")
	}
}

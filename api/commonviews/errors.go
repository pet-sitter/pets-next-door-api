package commonviews

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

func ParseBody(w http.ResponseWriter, r *http.Request, payload interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		BadRequest(w, nil, err.Error())
		return err
	}
	if err := validator.New().Struct(payload); err != nil {
		BadRequest(w, nil, err.Error())
		return err
	}

	return nil
}

type ErrorView struct {
	Message string `json:"message"`
}

func NewErrorView(message string) *ErrorView {
	return &ErrorView{
		Message: message,
	}
}

func BadRequest(w http.ResponseWriter, headers map[string]string, message string) error {
	if message == "" {
		message = "bad request"
	}

	return writePayload(w, headers, NewErrorView(message), http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter, headers map[string]string, message string) error {
	if message == "" {
		message = "unauthorized"
	}

	return writePayload(w, headers, NewErrorView(message), http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter, headers map[string]string, message string) error {
	if message == "" {
		message = "forbidden"
	}

	return writePayload(w, headers, NewErrorView(message), http.StatusForbidden)
}

func NotFound(w http.ResponseWriter, headers map[string]string, message string) error {
	if message == "" {
		message = "not found"
	}

	return writePayload(w, headers, NewErrorView(message), http.StatusNotFound)
}

func Conflict(w http.ResponseWriter, headers map[string]string, message string) error {
	if message == "" {
		message = "conflict"
	}

	return writePayload(w, headers, NewErrorView(message), http.StatusConflict)
}

func UnprocessableEntity(w http.ResponseWriter, headers map[string]string, message string) error {
	if message == "" {
		message = "unprocessable entity"
	}

	return writePayload(w, headers, NewErrorView(message), http.StatusUnprocessableEntity)
}

func InternalServerError(w http.ResponseWriter, headers map[string]string, message string) error {
	if message == "" {
		message = "internal server error"
	}

	return writePayload(w, headers, NewErrorView(message), http.StatusInternalServerError)
}

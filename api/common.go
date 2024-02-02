package pnd

import (
	"encoding/json"
	"net/http"
)

type PaginatedView[T interface{}] struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Items []T `json:"items"`
}

func NewPaginatedView[T interface{}](page int, size int, items []T) *PaginatedView[T] {
	return &PaginatedView[T]{
		Page:  page,
		Size:  size,
		Items: items,
	}
}

func writePayload(w http.ResponseWriter, headers map[string]string, payload interface{}, statusCode int) error {
	setHeaders(w, headers)

	w.WriteHeader(statusCode)

	if payload == nil {
		return nil
	}

	return json.NewEncoder(w).Encode(payload)
}

func setHeaders(w http.ResponseWriter, headers map[string]string) {
	if headers == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}

	for key, value := range headers {
		w.Header().Set(key, value)
	}
}

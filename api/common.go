package pnd

import (
	"encoding/json"
	"net/http"
)

type PaginatedView[T interface{}] struct {
	Page       int  `json:"page"`
	Size       int  `json:"size"`
	IsLastPage bool `json:"isLastPage"`
	Items      []T  `json:"items"`
}

func NewPaginatedView[T interface{}](page, size int, isLastPage bool, items []T) *PaginatedView[T] {
	return &PaginatedView[T]{
		Page:       page,
		Size:       size,
		IsLastPage: isLastPage,
		Items:      items,
	}
}

// CalcLastPage는 현재 페이지가 마지막 페이지인지를 items의 개수를 통해 계산한다.
// 마지막 페이지가 아니라면 items를 size만큼 잘라서 마지막 페이지로 만든다.
func (l *PaginatedView[T]) CalcLastPage() {
	if len(l.Items) > l.Size {
		l.IsLastPage = false
		l.Items = l.Items[:l.Size]
	} else {
		l.IsLastPage = true
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

type CursorPaginatedView[T interface{}] struct {
	Prev    string `json:"prev"`
	Next    string `json:"next"`
	HasPrev bool   `json:"hasPrev"`
	HasNext bool   `json:"hasNext"`
	Items   []T    `json:"items"`
}

func NewCursorPaginatedView[T interface{}](prev, next string, hasPrev, hasNext bool, items []T) *CursorPaginatedView[T] {
	return &CursorPaginatedView[T]{
		Prev:    prev,
		Next:    next,
		HasPrev: hasPrev,
		HasNext: hasNext,
		Items:   items,
	}
}

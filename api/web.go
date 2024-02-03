package pnd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

func ParseBody(r *http.Request, payload interface{}) *AppError {
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return ErrInvalidBody(err)
	}
	if err := validator.New().Struct(payload); err != nil {
		return ErrInvalidBody(err)
	}

	return nil
}

func ParseIdFromPath(r *http.Request, path string) (*int, *AppError) {
	id, err := strconv.Atoi(chi.URLParam(r, path))
	if err != nil {
		return nil, ErrInvalidParam(err)
	}
	if id <= 0 {
		return nil, ErrInvalidParam(fmt.Errorf("expected integer value bigger than 0 for path: %s", path))
	}

	return &id, nil
}

func ParseOptionalIntQuery(r *http.Request, query string) (*int, *AppError) {
	queryStr := r.URL.Query().Get(query)
	if queryStr == "" {
		return nil, nil
	}

	value, err := strconv.Atoi(queryStr)
	if err != nil {
		return nil, ErrInvalidQuery(fmt.Errorf("expected integer value for query: %s", query))
	}

	return &value, nil
}

func ParseRequiredStringQuery(r *http.Request, query string) (*string, *AppError) {
	queryStr := r.URL.Query().Get(query)
	if queryStr == "" {
		return nil, ErrInvalidQuery(fmt.Errorf("expected non-empty string for query: %s", query))
	}

	return &queryStr, nil
}

func ParseOptionalStringQuery(r *http.Request, query string) *string {
	queryStr := r.URL.Query().Get(query)
	if queryStr == "" {
		return nil
	}

	return &queryStr
}

// ParsePaginationQueries parses pagination parameters from query string: page, size.
func ParsePaginationQueries(r *http.Request, defaultPage int, defaultLimit int) (page int, size int, err *AppError) {
	pageQuery := r.URL.Query().Get("page")
	sizeQuery := r.URL.Query().Get("size")

	page = defaultPage
	size = defaultLimit

	if pageQuery != "" {
		var atoiError error
		page, atoiError = strconv.Atoi(pageQuery)
		if atoiError != nil {
			return 0, 0, ErrInvalidPagination(fmt.Errorf("expected integer value bigger than 0 for query: page"))
		}
	}

	if sizeQuery != "" {
		var atoiError error
		size, atoiError = strconv.Atoi(sizeQuery)
		if atoiError != nil {
			return 0, 0, ErrInvalidPagination(fmt.Errorf("expected integer value bigger than 0 for query: size"))
		}
	}

	return page, size, nil
}

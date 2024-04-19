package pnd

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ParseBody(c echo.Context, payload interface{}) *AppError {
	if err := c.Bind(payload); err != nil {
		return ErrInvalidBody(err)
	}
	if err := validator.New().Struct(payload); err != nil {
		return ErrInvalidBody(err)
	}

	return nil
}

func ParseIDFromPath(c echo.Context, path string) (*int, *AppError) {
	id, err := strconv.Atoi(c.Param(path))
	if err != nil {
		return nil, ErrInvalidParam(err)
	}
	if id <= 0 {
		return nil, ErrInvalidParam(fmt.Errorf("expected integer value bigger than 0 for path: %s", path))
	}

	return &id, nil
}

func ParseOptionalIntQuery(c echo.Context, query string) (*int, *AppError) {
	queryStr := c.QueryParam(query)
	if queryStr == "" {
		return nil, nil
	}

	value, err := strconv.Atoi(queryStr)
	if err != nil {
		return nil, ErrInvalidQuery(fmt.Errorf("expected integer value for query: %s", query))
	}

	return &value, nil
}

func ParseRequiredStringQuery(c echo.Context, query string) (*string, *AppError) {
	queryStr := c.QueryParam(query)
	if queryStr == "" {
		return nil, ErrInvalidQuery(fmt.Errorf("expected non-empty string for query: %s", query))
	}

	return &queryStr, nil
}

func ParseOptionalStringQuery(c echo.Context, query string) *string {
	queryStr := c.QueryParam(query)
	if queryStr == "" {
		return nil
	}

	return &queryStr
}

// ParsePaginationQueries parses pagination parameters from query string: page, size.
func ParsePaginationQueries(c echo.Context, defaultPage, defaultLimit int) (page, size int, err *AppError) {
	pageQuery := c.QueryParam("page")
	sizeQuery := c.QueryParam("size")

	page = defaultPage
	size = defaultLimit

	if pageQuery != "" {
		var atoiError error
		page, atoiError = strconv.Atoi(pageQuery)
		if atoiError != nil || page <= 0 {
			return 0, 0, ErrInvalidPagination(fmt.Errorf("expected integer value bigger than 0 for query: page"))
		}
	}

	if sizeQuery != "" {
		var atoiError error
		size, atoiError = strconv.Atoi(sizeQuery)
		if atoiError != nil || size <= 0 {
			return 0, 0, ErrInvalidPagination(fmt.Errorf("expected integer value bigger than 0 for query: size"))
		}
	}

	return page, size, nil
}

package pnd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ParseBody(c echo.Context, payload interface{}) error {
	if err := c.Bind(payload); err != nil {
		return ErrInvalidBody(err)
	}
	if err := validator.New().Struct(payload); err != nil {
		return ErrInvalidBody(err)
	}

	return nil
}

func ParseIDFromPath(c echo.Context, path string) (uuid.UUID, error) {
	idStr := c.Param(path)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, ErrInvalidParam(fmt.Errorf("expected valid UUID for path: %s", path))
	}

	return id, nil
}

func ParseOptionalUUIDQuery(c echo.Context, query string) (uuid.NullUUID, error) {
	queryStr := c.QueryParam(query)
	if queryStr == "" {
		return uuid.NullUUID{}, nil
	}

	id, err := uuid.Parse(queryStr)
	if err != nil {
		return uuid.NullUUID{}, ErrInvalidQuery(
			fmt.Errorf("expected valid UUID for query: %s", query),
		)
	}

	return uuid.NullUUID{UUID: id, Valid: true}, nil
}

func ParseOptionalIntQuery(c echo.Context, query string) (*int, error) {
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

func ParseRequiredStringQuery(c echo.Context, query string) (*string, error) {
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
func ParsePaginationQueries(
	c echo.Context,
	defaultPage, defaultLimit int,
) (page, size int, err error) {
	pageQuery := c.QueryParam("page")
	sizeQuery := c.QueryParam("size")

	page = defaultPage
	size = defaultLimit

	if pageQuery != "" {
		var atoiError error
		page, atoiError = strconv.Atoi(pageQuery)
		if atoiError != nil || page <= 0 {
			return 0, 0, ErrInvalidPagination(
				errors.New("expected integer value bigger than 0 for query: page"),
			)
		}
	}

	if sizeQuery != "" {
		var atoiError error
		size, atoiError = strconv.Atoi(sizeQuery)
		if atoiError != nil || size <= 0 {
			return 0, 0, ErrInvalidPagination(
				errors.New("expected integer value bigger than 0 for query: size"),
			)
		}
	}

	return page, size, nil
}

func ParseCursorPaginationQueries(
	c echo.Context, defaultLimit int,
) (prev, next uuid.NullUUID, limit int, err error) {
	prevQuery := c.QueryParam("prev")
	nextQuery := c.QueryParam("next")
	sizeQuery := c.QueryParam("size")

	if prevQuery != "" {
		id, err := uuid.Parse(prevQuery)
		if err != nil {
			return uuid.NullUUID{},
				uuid.NullUUID{},
				0,
				ErrInvalidQuery(errors.New("expected valid UUID for query: prev"))
		}
		prev = uuid.NullUUID{UUID: id, Valid: true}
	}

	if nextQuery != "" {
		id, err := uuid.Parse(nextQuery)
		if err != nil {
			return uuid.NullUUID{},
				uuid.NullUUID{},
				0,
				ErrInvalidQuery(errors.New("expected valid UUID for query: next"))
		}
		next = uuid.NullUUID{UUID: id, Valid: true}
	}

	if sizeQuery != "" {
		var atoiError error
		limit, atoiError = strconv.Atoi(sizeQuery)
		if atoiError != nil || limit <= 0 {
			return uuid.NullUUID{},
				uuid.NullUUID{},
				0,
				ErrInvalidQuery(errors.New("expected integer value bigger than 0 for query: size"))
		}
	}

	if limit == 0 {
		limit = defaultLimit
	}

	return prev, next, limit, nil
}

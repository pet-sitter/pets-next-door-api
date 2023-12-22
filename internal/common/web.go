package webutils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseIdFromPath(r *http.Request, path string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, path))
}

// ParsePaginationQueries parses pagination parameters from query string: page, size.
func ParsePaginationQueries(r *http.Request, defaultPage int, defaultLimit int) (page int, size int, err error) {
	pageQuery := r.URL.Query().Get("page")
	sizeQuery := r.URL.Query().Get("size")

	page = defaultPage
	size = defaultLimit

	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			return 0, 0, err
		}
	}

	if sizeQuery != "" {
		size, err = strconv.Atoi(sizeQuery)
		if err != nil {
			return 0, 0, err
		}
	}

	return page, size, nil
}

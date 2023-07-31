package webutils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func ParseIdFromPath(r *http.Request, path string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, path))
}

package views

import "net/http"

func OK(w http.ResponseWriter, headers map[string]string, payload interface{}) error {
	return writePayload(w, headers, payload, http.StatusOK)
}

func Created(w http.ResponseWriter, headers map[string]string, payload interface{}) error {
	return writePayload(w, headers, payload, http.StatusCreated)
}

func NoContent(w http.ResponseWriter, headers map[string]string) error {
	return writePayload(w, headers, nil, http.StatusNoContent)
}

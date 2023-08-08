package commonviews

import (
	"encoding/json"
	"net/http"
)

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
		w.Header().Set("Content-Type", "application/json")
		return
	}

	for key, value := range headers {
		w.Header().Set(key, value)
	}
}

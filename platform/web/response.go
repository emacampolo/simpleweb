package web

import (
	"encoding/json"
	"net/http"
)

// EncodeJSON encodes the given data as JSON and writes it to the ResponseWriter.
func EncodeJSON(w http.ResponseWriter, data any, statusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	var jsonData []byte

	var err error
	switch v := data.(type) {
	case []byte:
		jsonData = v
	default:
		jsonData, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

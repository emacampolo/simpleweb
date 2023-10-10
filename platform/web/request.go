package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Param returns the value of the URL parameter with the given key.
// If the parameter is not found, it returns an empty string.
func Param(r *http.Request, key string) string {
	return chi.RouteContext(r.Context()).URLParam(key)
}

// ParamInt returns the value of the URL parameter with the given key as an int.
// If the parameter is not found, it returns 0.
// If the parameter type value is a not an int, it returns an error.
func ParamInt(r *http.Request, key string) (int, error) {
	value := Param(r, key)
	if value == "" {
		return 0, nil
	}

	intValue, err := strconv.Atoi(value)
	return intValue, err
}

// WithURLParams adds the given URL parameters to the request context.
// testing.T is required but not used to enforce the use of this function in tests only.
func WithURLParams(t *testing.T, req *http.Request, params map[string]string) *http.Request {
	if t == nil {
		panic("use WithURLParams only in tests")
	}

	var routeParams chi.RouteParams
	for key, val := range params {
		routeParams.Add(key, val)
	}

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams = routeParams

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	return req
}

var _validator = validator.New()

// DecodeJSON decodes the request body into the provided value.
func DecodeJSON(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return _validator.Struct(val)
}

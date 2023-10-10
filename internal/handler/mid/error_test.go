package mid_test

import (
	"bootcamp-web/internal/handler"
	"bootcamp-web/internal/handler/mid"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewErrorReturnNilIfHandlerReturnNil(t *testing.T) {
	errMiddleware := mid.NewError()
	err := errMiddleware(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})(nil, nil)

	require.NoError(t, err)
}

func TestNewErrorReturnInternalServerErrorIfHandlerAnyError(t *testing.T) {
	errMiddleware := mid.NewError()

	request := httptest.NewRequest("GET", "/", nil)
	writer := httptest.NewRecorder()

	err := errMiddleware(func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("any error")
	})(writer, request)

	require.NoError(t, err)
	response := writer.Result()
	require.Equal(t, 500, response.StatusCode)
	require.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	require.Equal(t, `{"code":"internal_server_error","message":"any error"}`, writer.Body.String())
}

func TestNewErrorReturnCustomError(t *testing.T) {
	errMiddleware := mid.NewError()

	request := httptest.NewRequest("GET", "/", nil)
	writer := httptest.NewRecorder()

	err := errMiddleware(func(w http.ResponseWriter, r *http.Request) error {
		return handler.NewError(http.StatusBadRequest, "invalid request")
	})(writer, request)

	require.NoError(t, err)
	response := writer.Result()
	require.Equal(t, 400, response.StatusCode)
	require.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	require.Equal(t, `{"code":"bad_request","message":"invalid request"}`, writer.Body.String())
}

package web_test

import (
	"bootcamp-web/platform/web"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"hello": "world"}
	err := web.EncodeJSON(w, data, 200)

	require.NoError(t, err)
	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.Equal(t, 200, w.Code)
	require.Equal(t, "{\"hello\":\"world\"}", w.Body.String())
}

func TestEncodeJSON_Error(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)
	err := web.EncodeJSON(w, data, 200)

	require.EqualError(t, err, "json: unsupported type: chan int")
}

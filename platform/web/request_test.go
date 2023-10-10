package web_test

import (
	"bootcamp-web/platform/web"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParam(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = web.WithURLParams(t, r, map[string]string{"id": "123"})
	require.Equal(t, "123", web.Param(r, "id"))
}

func TestWithURLParams_PanicIfTestingTIsNil(t *testing.T) {
	require.PanicsWithValue(t, "use WithURLParams only in tests", func() {
		web.WithURLParams(nil, httptest.NewRequest("GET", "/", nil), nil)
	})
}

func TestDecodeJSON(t *testing.T) {
	r := httptest.NewRequest("GET", "/", strings.NewReader(`{"name":"bill"}`))
	var val struct {
		Name string `json:"name"`
	}

	err := web.DecodeJSON(r, &val)
	require.NoError(t, err)
	require.Equal(t, "bill", val.Name)
}

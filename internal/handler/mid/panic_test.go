package mid_test

import (
	"bootcamp-web/internal/handler/mid"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPanic(t *testing.T) {
	panicMiddleware := mid.NewPanic()
	err := panicMiddleware(func(w http.ResponseWriter, r *http.Request) error {
		panic("any panic")
	})(nil, nil)

	require.ErrorContains(t, err, "panic [any panic] trace [goroutine")
}

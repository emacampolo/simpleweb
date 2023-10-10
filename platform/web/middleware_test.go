package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_wrapMiddlewareIterateBackwards(t *testing.T) {
	var names []string
	midFactory := func(name string) Middleware {
		return func(next Handler) Handler {
			return func(w http.ResponseWriter, r *http.Request) error {
				names = append(names, name)
				return next(w, r)
			}
		}
	}

	middlewares := []Middleware{
		midFactory("1"),
		midFactory("2"),
		midFactory("3"),
	}

	wrapped := wrapMiddleware(middlewares, func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})

	err := wrapped(nil, nil)
	require.Nil(t, err)
	require.Equal(t, names, []string{"1", "2", "3"})
}

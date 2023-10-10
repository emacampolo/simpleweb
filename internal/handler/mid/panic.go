package mid

import (
	"bootcamp-web/platform/web"
	"fmt"
	"net/http"
	"runtime/debug"
)

// NewPanic handles any panic that may occur by notifying the error to an external system
// and responding to the client with status code 500.
func NewPanic() web.Middleware {
	return func(webHandler web.Handler) web.Handler {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("panic [%v] trace [%s]", rec, string(trace))
				}
			}()

			return webHandler(w, r)
		}
	}
}

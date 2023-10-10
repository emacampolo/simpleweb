package mid

import (
	"bootcamp-web/internal/handler"
	"bootcamp-web/platform/web"
	"encoding/json"
	"errors"
	"net/http"
)

// NewError is a middleware that handles errors returned by handlers.
func NewError() web.Middleware {
	return func(webHandler web.Handler) web.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			err := webHandler(w, r)
			if err == nil {
				return nil
			}

			var webErr *handler.Error
			if !errors.As(err, &webErr) {
				err = handler.NewError(http.StatusInternalServerError, err.Error())
			}

			contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
			if m, ok := err.(json.Marshaler); ok {
				if jsonBody, marshalErr := m.MarshalJSON(); marshalErr == nil {
					contentType, body = "application/json; charset=utf-8", jsonBody
				}
			}

			w.Header().Set("Content-Type", contentType)
			if h, ok := err.(interface{ Headers() http.Header }); ok {
				for k, values := range h.Headers() {
					for _, v := range values {
						w.Header().Add(k, v)
					}
				}
			}

			code := http.StatusInternalServerError
			if sc, ok := err.(interface{ StatusCode() int }); ok {
				code = sc.StatusCode()
			}

			w.WriteHeader(code)
			_, _ = w.Write(body)

			// If any error is returned out of the middleware stack, the application will
			// panic. So, we return nil here to indicate that the error has been
			// handled.
			return nil
		}
	}
}

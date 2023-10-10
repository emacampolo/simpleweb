package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler is the signature used by all application handlers.
// The error returned is intended to be used be the error handling middleware.
// See internal/handler/mid/error.go.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Muxer is a simple HTTP route multiplexer.
type Muxer struct {
	mux *chi.Mux
	mw  []Middleware
}

// NewMux returns a new Muxer instance.
// The shutdown channel is used to gracefully shut down the app when an integrity issue is identified.
func NewMux(mw ...Middleware) *Muxer {
	return &Muxer{
		mux: chi.NewRouter(),
		mw:  mw,
	}
}

// HandleNoMiddleware associates a handler function with an HTTP method and URL pattern.
// It does not include the application middleware.
func (s *Muxer) HandleNoMiddleware(method string, path string, handler Handler) {
	s.handle(method, path, handler)
}

// Handle associates a handler function with an HTTP method and URL pattern.
// The mw list should contain middleware (or middleware chains) to be applied.
func (s *Muxer) Handle(method string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(s.mw, handler)

	s.handle(method, path, handler)
}

func (s *Muxer) handle(method string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {

			// If we get to this point it means that neither the handler nor the
			// middleware stack could handle the error. So, we just panic.
			panic(err)
		}
	}

	s.mux.MethodFunc(method, path, h)
}

// ServeHTTP implements the http.Handler interface.
func (s *Muxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

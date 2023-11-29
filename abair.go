package abair

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-chi/chi/v5"
)

var loadErrorHandler sync.Once

// Server is a wrapper around chi.Router.
type Server struct {
	Router       chi.Router
	ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

// ServeHTTP implements http.Handler.
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// Request is a request.
type Request[T any] struct {
	Body    T
	Params  url.Values
	Headers http.Header
}

// HandlerFunc is a handler function.
type HandlerFunc[Req, Resp any] func(context.Context, Request[Req]) (Resp, error)

// Get is a GET handler.
func Get[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Get(path, handler(s, hndlr))
}

// Post is a POST handler.
func Post[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Post(path, handler(s, hndlr))
}

// Put is a PUT handler.
func Put[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Put(path, handler(s, hndlr))
}

// Delete is a DELETE handler.
func Delete[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Delete(path, handler(s, hndlr))
}

// Patch is a PATCH handler.
func Patch[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Patch(path, handler(s, hndlr))
}

// Options is a OPTIONS handler.
func Options[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Options(path, handler(s, hndlr))
}

// Head is a HEAD handler.
func Head[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Head(path, handler(s, hndlr))
}

// Connect is a CONNECT handler.
func Connect[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Connect(path, handler(s, hndlr))
}

// Trace is a TRACE handler.
func Trace[Req, Resp any](s *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	s.Router.Trace(path, handler(s, hndlr))
}

func handler[Req, Resp any](s *Server, hndlr HandlerFunc[Req, Resp]) http.HandlerFunc {
	loadErrorHandler.Do(func() {
		if s.ErrorHandler == nil {
			s.ErrorHandler = defaultErrorHandler
		}
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req Request[Req]
		if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil {
			s.ErrorHandler(w, r, err)
			return
		}

		req.Params = r.URL.Query()
		req.Headers = r.Header

		resp, err := hndlr(ctx, req)
		if err != nil {
			s.ErrorHandler(w, r, err)
			return
		}

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			s.ErrorHandler(w, r, err)
			return
		}
	})
}

func defaultErrorHandler(w http.ResponseWriter, _ *http.Request, _ error) {
	w.WriteHeader(http.StatusInternalServerError)
}

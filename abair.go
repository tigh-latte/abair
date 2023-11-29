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
func (f Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.Router.ServeHTTP(w, r)
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
func Get[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Get(path, handler(f, hndlr))
}

// Post is a POST handler.
func Post[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Post(path, handler(f, hndlr))
}

// Put is a PUT handler.
func Put[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Put(path, handler(f, hndlr))
}

// Delete is a DELETE handler.
func Delete[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Delete(path, handler(f, hndlr))
}

// Patch is a PATCH handler.
func Patch[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Patch(path, handler(f, hndlr))
}

// Options is a OPTIONS handler.
func Options[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Options(path, handler(f, hndlr))
}

// Head is a HEAD handler.
func Head[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Head(path, handler(f, hndlr))
}

// Connect is a CONNECT handler.
func Connect[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Connect(path, handler(f, hndlr))
}

// Trace is a TRACE handler.
func Trace[Req, Resp any](f *Server, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Trace(path, handler(f, hndlr))
}

func handler[Req, Resp any](f *Server, hndlr HandlerFunc[Req, Resp]) http.HandlerFunc {
	loadErrorHandler.Do(func() {
		if f.ErrorHandler == nil {
			f.ErrorHandler = defaultErrorHandler
		}
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req Request[Req]
		if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil {
			f.ErrorHandler(w, r, err)
			return
		}

		req.Params = r.URL.Query()
		req.Headers = r.Header

		resp, err := hndlr(ctx, req)
		if err != nil {
			f.ErrorHandler(w, r, err)
			return
		}

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			f.ErrorHandler(w, r, err)
			return
		}
	})
}

func defaultErrorHandler(w http.ResponseWriter, _ *http.Request, _ error) {
	w.WriteHeader(http.StatusInternalServerError)
}

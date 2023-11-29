package abair

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

// Framework is a wrapper around chi.Router.
type Framework struct {
	Router       chi.Router
	ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

func (f Framework) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.Router.ServeHTTP(w, r)
}

// Request is a request.
type Request[T any] struct {
	Body   T
	Params url.Values
	Header http.Header
}

// HandlerFunc is a handler function.
type HandlerFunc[Req, Resp any] func(context.Context, Request[Req]) (Resp, error)

func Get[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Get(path, handler(f, path, hndlr))
}

func Post[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Post(path, handler(f, path, hndlr))
}

func Put[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Put(path, handler(f, path, hndlr))
}

func Delete[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Delete(path, handler(f, path, hndlr))
}

func Patch[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Patch(path, handler(f, path, hndlr))
}

func Options[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Options(path, handler(f, path, hndlr))
}

func Head[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Head(path, handler(f, path, hndlr))
}

func Connect[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Connect(path, handler(f, path, hndlr))
}

func Trace[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) {
	f.Router.Trace(path, handler(f, path, hndlr))
}

func handler[Req, Resp any](f *Framework, path string, hndlr HandlerFunc[Req, Resp]) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req Request[Req]
		if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil {
			f.ErrorHandler(w, r, err)
			return
		}

		resp, err := hndlr(ctx, req)
		if err != nil {
			f.ErrorHandler(w, r, err)
			return
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			f.ErrorHandler(w, r, err)
			return
		}
	})
}

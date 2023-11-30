package abair

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-chi/chi/v5"
)

var loadErrorHandler sync.Once

func applyDefaultServerCfg(s *Server) func() {
	return func() {
		if s.Logger == nil {
			s.Logger = slog.Default()
		}
		if s.ErrorHandler == nil {
			s.ErrorHandler = buildDefaultErrorHandler(s.Logger)
		}
	}
}

// Server is a wrapper around chi.Router.
type Server struct {
	Router       chi.Router
	Logger       *slog.Logger
	ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

// ServeHTTP implements http.Handler.
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// Request is a request.
type Request[T any] struct {
	Body        T
	PathParams  map[string]string
	QueryParams url.Values
	Headers     http.Header
}

// HandlerFunc is a handler function.
type HandlerFunc[Req, Resp any] func(context.Context, Request[Req]) (Resp, error)

// Route is a route.
func Route(s *Server, path string, fn func(s *Server)) {
	loadErrorHandler.Do(applyDefaultServerCfg(s))
	sub := &Server{
		Router:       chi.NewRouter(),
		Logger:       s.Logger,
		ErrorHandler: s.ErrorHandler,
	}
	fn(sub)
	s.Router.Mount(path, sub)
}

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
	loadErrorHandler.Do(applyDefaultServerCfg(s))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		routeCtx := chi.RouteContext(ctx)

		req := Request[Req]{
			PathParams:  make(map[string]string, len(routeCtx.URLParams.Keys)),
			Headers:     r.Header,
			Body:        *new(Req),
			QueryParams: r.URL.Query(),
		}
		if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil && err != io.EOF {
			s.ErrorHandler(w, r, err)
			return
		}

		for i := 0; i < len(routeCtx.URLParams.Keys); i++ {
			req.PathParams[routeCtx.URLParams.Keys[i]] = routeCtx.URLParams.Values[i]
		}

		req.QueryParams = r.URL.Query()
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

func buildDefaultErrorHandler(log *slog.Logger) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h := &HTTPError{}
		if ok := errors.As(err, &h); ok {
			if h.Internal != nil {
				if herr, ok := h.Internal.(*HTTPError); ok {
					h = herr
				}
			}
		} else {
			var h2 HTTPError
			fmt.Println(errors.As(err, &h2))
			h = &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}

		var (
			code    = h.Code
			message = h.Message
		)

		response := map[string]any{}

		switch m := message.(type) {
		case string:
			response["message"] = m
		case json.Marshaler:
			// do nothing
		case error:
			response["message"] = m.Error()
		}

		w.WriteHeader(code)
		if r.Method == http.MethodHead {
			return
		}

		bb, err := json.Marshal(response)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError,
				"failed to encode error",
				slog.Any("error", err),
			)
		}

		_, _ = w.Write(bb)
	}
}

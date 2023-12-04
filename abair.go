package abair

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v3"
)

type Router = chi.Router

// Server is a wrapper around chi.Router.
type Server struct {
	Router
	Logger         *slog.Logger
	Decoders       map[string]func(io.Reader) Decoder
	Encoders       map[string]func(io.Writer) Encoder
	DefaultDecoder func(io.Reader) Decoder
	DefaultEncoder func(io.Writer) Encoder
	ErrorHandler   func(w http.ResponseWriter, r *http.Request, err error)

	PreferredResponseType string
}

// NewServer returns a new Server.
func NewServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
		Logger: slog.Default(),
		DefaultDecoder: func(r io.Reader) Decoder {
			return json.NewDecoder(r)
		},
		DefaultEncoder: func(w io.Writer) Encoder {
			return json.NewEncoder(w)
		},
		PreferredResponseType: "application/json",
	}
	s.ErrorHandler = buildDefaultErrorHandler(s, slog.Default())
	s.Decoders = map[string]func(io.Reader) Decoder{
		"application/json": func(r io.Reader) Decoder { return json.NewDecoder(r) },
		"application/xml":  func(r io.Reader) Decoder { return xml.NewDecoder(r) },
		"application/yaml": func(r io.Reader) Decoder { return yaml.NewDecoder(r) },
		"text/xml":         func(r io.Reader) Decoder { return xml.NewDecoder(r) },
		"text/yaml":        func(r io.Reader) Decoder { return yaml.NewDecoder(r) },
	}
	s.Encoders = map[string]func(io.Writer) Encoder{
		"application/json": func(w io.Writer) Encoder { return json.NewEncoder(w) },
		"application/xml":  func(w io.Writer) Encoder { return xml.NewEncoder(w) },
		"application/yaml": func(w io.Writer) Encoder { return yaml.NewEncoder(w) },
		"text/xml":         func(w io.Writer) Encoder { return xml.NewEncoder(w) },
		"text/yaml":        func(w io.Writer) Encoder { return yaml.NewEncoder(w) },
	}
	s.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.ErrorHandler(w, r, &HTTPError{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		})
	})
	s.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		s.ErrorHandler(w, r, &HTTPError{
			Code:    http.StatusMethodNotAllowed,
			Message: http.StatusText(http.StatusMethodNotAllowed),
		})
	})

	return s
}

// Request is a request.
type Request[B, P any] struct {
	Body        B
	PathParams  P
	QueryParams url.Values
	Headers     http.Header
}

// HandlerFunc is a handler function.
type HandlerFunc[Body, Path, Resp any] func(context.Context, Request[Body, Path]) (Resp, error)

// Route is a route.
func (s *Server) Route(path string, fn func(s *Server)) {
	sub := &Server{
		Router:                chi.NewRouter(),
		Logger:                s.Logger,
		DefaultDecoder:        s.DefaultDecoder,
		DefaultEncoder:        s.DefaultEncoder,
		Decoders:              s.Decoders,
		Encoders:              s.Encoders,
		PreferredResponseType: s.PreferredResponseType,
		ErrorHandler:          s.ErrorHandler,
	}
	fn(sub)
	s.Mount(path, sub)
}

// HTTPHandlerWrapper wraps a abair.HandlerFunc to be used as a http.HandlerFunc.
func HTTPHandlerWrapper[Body, Path, Resp any](s *Server, hndlr HandlerFunc[Body, Path, Resp]) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := Request[Body, Path]{
			PathParams:  *new(Path),
			Headers:     r.Header,
			Body:        *new(Body),
			QueryParams: r.URL.Query(),
		}
		if err := s.decode(r, &req.Body); err != nil && err != io.EOF {
			s.ErrorHandler(w, r, err)
			return
		}

		var (
			pathType = reflect.TypeOf(req.PathParams)
			pathVal  = reflect.ValueOf(&req.PathParams)
		)

		for i := 0; i < pathType.NumField(); i++ {
			field := pathType.Field(i)
			if field.PkgPath == "" {
				continue // skip unexported fields.
			}
			tag := field.Tag.Get("path")
			if tag == "" {
				tag = field.Name
			}

			pval := chi.URLParamFromCtx(ctx, tag)
			if pval == "" {
				s.ErrorHandler(w, r, &HTTPError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("missing param %q", tag),
				})
				return
			}

			switch field.Type.Kind() {
			case reflect.String:
				pathVal.Elem().Field(i).SetString(pval)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val, err := strconv.ParseInt(pval, 10, 64)
				if err != nil {
					s.ErrorHandler(w, r, &HTTPError{
						Code:     http.StatusBadRequest,
						Message:  fmt.Sprintf("expected param %q to be an integer", tag),
						Internal: err,
					})
					return
				}
				pathVal.Elem().Field(i).SetInt(int64(val))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val, err := strconv.ParseUint(pval, 10, 64)
				if err != nil {
					s.ErrorHandler(w, r, &HTTPError{
						Code:     http.StatusBadRequest,
						Internal: err,
					})
					return
				}
				pathVal.Elem().Field(i).SetUint(val)
			case reflect.Float32, reflect.Float64:
				val, err := strconv.ParseFloat(pval, 64)
				if err != nil {
					s.ErrorHandler(w, r, &HTTPError{
						Code:     http.StatusBadRequest,
						Internal: err,
					})
					return
				}
				pathVal.Elem().Field(i).SetFloat(val)
			case reflect.Struct:
				loader, ok := pathVal.Elem().Field(i).Addr().Interface().(interface {
					ParsePath(string) error
				})
				if !ok {
					break
				}
				if err := loader.ParsePath(pval); err != nil {
					s.ErrorHandler(w, r, err)
					return
				}
			case reflect.Ptr:
				item := pathVal.Elem().Field(i)
				if item.IsNil() {
					item.Set(reflect.New(field.Type.Elem()))
				}
				loader, ok := item.Interface().(interface {
					ParsePath(string) error
				})
				if !ok {
					break
				}
				if err := loader.ParsePath(pval); err != nil {
					s.ErrorHandler(w, r, err)
					return
				}
			}
		}

		resp, err := hndlr(ctx, req)
		if err != nil {
			s.ErrorHandler(w, r, err)
			return
		}

		if err = s.encode(w, r, resp); err != nil {
			s.ErrorHandler(w, r, err)
			return
		}
	})
}

func buildDefaultErrorHandler(s *Server, log *slog.Logger) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h := &HTTPError{}
		if ok := errors.As(err, &h); ok {
			if h.Internal != nil {
				if herr, ok := h.Internal.(*HTTPError); ok {
					h = herr
				}
			}
		} else {
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
		case json.Marshaler, xml.Marshaler, yaml.Marshaler:
			// do nothing
		case error:
			response["message"] = m.Error()
		}

		w.WriteHeader(code)
		if r.Method == http.MethodHead {
			return
		}

		if err := s.encode(w, r, response); err != nil {
			log.LogAttrs(
				r.Context(),
				slog.LevelError,
				"failed to encode error",
				slog.Any("error", err),
			)
		}
	}
}

func (s *Server) decode(r *http.Request, v any) error {
	fn, ok := s.Decoders[r.Header.Get("Content-Type")]
	if !ok {
		return s.DefaultDecoder(r.Body).Decode(v)
	}
	return fn(r.Body).Decode(v)
}

func (s *Server) encode(w http.ResponseWriter, r *http.Request, v any) error {
	accept := r.Header.Get("Accept")
	if idx := strings.IndexByte(accept, ','); idx >= 0 {
		if strings.Contains(accept, s.PreferredResponseType) {
			accept = s.PreferredResponseType
		} else {
			accept = accept[:idx]
		}
	}
	if idx := strings.IndexByte(accept, ';'); idx >= 0 {
		accept = accept[:idx]
	}
	fn, ok := s.Encoders[accept]
	if !ok {
		w.Header().Set("Content-Type", s.PreferredResponseType)
		return s.DefaultEncoder(w).Encode(v)
	}

	w.Header().Set("Content-Type", accept)
	return fn(w).Encode(v)
}

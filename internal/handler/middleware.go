package handler

import (
	// "bytes"
	"context"
	// "io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type contextKey string

const userCtxKey = contextKey("role")

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (s *Server) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"message":"missing token"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		if strings.HasPrefix(token, s.Cfg.Auth.DummyTokenPrefix) {
			role := strings.TrimPrefix(token, s.Cfg.Auth.DummyTokenPrefix)
			ctx := context.WithValue(r.Context(), userCtxKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims, err := s.JWTManager.Parse(token)
		if err != nil {
			http.Error(w, `{"message":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxRole, _ := r.Context().Value(userCtxKey).(string)
			if ctxRole != role {
				http.Error(w, `{"message":"forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// var bodyBytes []byte
		// if r.Body != nil {
		// 	bodyBytes, _ = io.ReadAll(r.Body)
		// }
		// r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		ww := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		slog.Info("HTTP Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", ww.statusCode),
			//slog.String("body", string(bodyBytes)),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	Status       int
	ResponseBody string
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *StatusRecorder) Write(body []byte) (int, error) {
	r.ResponseBody = string(body)
	return r.ResponseWriter.Write(body)
}

func (s *Server) PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		path := r.URL.Path // fallback
		if rc := chi.RouteContext(r.Context()); rc != nil {
			if p := rc.RoutePattern(); p != "" {
				path = p
			}
		}

		// request, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	slog.Info("http_logger: read request body")
		// }
		// r.Body = io.NopCloser(bytes.NewBuffer(request))

		start := time.Now()
		next.ServeHTTP(recorder, r)
		s.metrics.SaveHTTPDuration(start, path, recorder.Status, r.Method)
	})
}

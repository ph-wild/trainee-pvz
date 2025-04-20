package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type contextKey string

const userCtxKey = contextKey("role")

type statusRecorder struct {
	http.ResponseWriter
	Status       int
	ResponseBody string
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	r.ResponseBody = string(body)
	return r.ResponseWriter.Write(body)
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

		ww := &statusRecorder{ResponseWriter: w, Status: http.StatusOK}
		next.ServeHTTP(ww, r)

		slog.Info("HTTP Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", ww.Status),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func (s *Server) PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		path := r.URL.Path // fallback
		if rc := chi.RouteContext(r.Context()); rc != nil {
			if p := rc.RoutePattern(); p != "" {
				path = p
			}
		}

		start := time.Now()
		next.ServeHTTP(recorder, r)
		s.metrics.SaveHTTPDuration(start, path, recorder.Status, r.Method)
	})
}

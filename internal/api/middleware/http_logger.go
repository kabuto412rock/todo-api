package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// HTTPLogger is a top-level chi middleware.
// Logs request metadata, optional JSON request body, status, latency and JSON response body.
// Place early (r.Use) to reliably capture raw request & response.
func HTTPLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Request ID (reuse header or generate a new one)
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				id := uuid.New()
				requestID = base64.RawURLEncoding.EncodeToString(id[:])
			}
			// Read JSON request body (limit 1MB) & restore for handlers
			var reqBodyPreview string
			if allowBody(r.Method) && isJSON(r.Header.Get("Content-Type")) && r.Body != nil {
				raw, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB cap
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewBuffer(raw))
				if len(raw) > 0 {
					var js string
					if err := json.Unmarshal(raw, &js); err == nil {
						reqBodyPreview = js
					} else {
						snippet := string(raw)
						if len(snippet) > 512 {
							snippet = snippet[:512] + "... (truncated)"
						}
						reqBodyPreview = snippet
					}
				}
			}

			rec := &respRecorder{ResponseWriter: w, status: 200}

			attrs := []any{"id", requestID, "method", r.Method, "path", r.URL.Path}
			if r.URL.RawQuery != "" {
				attrs = append(attrs, "query", r.URL.RawQuery)
			}
			if len(reqBodyPreview) > 0 {
				attrs = append(attrs, slog.String("json_body", reqBodyPreview))
			}
			slog.Info("http request", attrs...)

			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			endAttrs := []any{"id", requestID, "status", rec.status, "duration_ms", duration.Milliseconds()}
			if len(rec.body) > 0 && isJSON(rec.Header().Get("Content-Type")) {
				snippet := rec.body
				if len(snippet) > 1024 {
					truncated := make([]byte, 1024)
					copy(truncated, snippet[:1024])
					truncated = append(truncated, []byte("... (truncated)")...)
					snippet = truncated
				}
				endAttrs = append(endAttrs, "response_body", string(snippet))
			}
			slog.Info("http request", endAttrs...)
		})
	}
}

type respRecorder struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (r *respRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *respRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

func allowBody(m string) bool {
	switch strings.ToUpper(m) {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// isJSON: minimal content-type check
func isJSON(ct string) bool {
	ct = strings.ToLower(ct)
	return strings.Contains(ct, "application/json") || strings.Contains(ct, "+json")
}

package handlers

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var TotalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "my_http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var ResponseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "my_response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

func (h *PostHandler) prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		ResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		TotalRequests.WithLabelValues(path).Inc()
	})
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version   = "dev"
	isReady   atomic.Bool
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200); w.Write([]byte("OK")) })
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		if isReady.Load() { w.WriteHeader(200); w.Write([]byte("READY")) } else { w.WriteHeader(503); w.Write([]byte("NOT_READY")) }
	})
	mux.HandleFunc("GET /api/v1/hello", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		hostname, _ := os.Hostname()
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello!", "hostname": hostname})
	}))
	mux.Handle("GET /metrics", promhttp.Handler())

	server := &http.Server{Addr: ":8080", Handler: mux}
	
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Starting server", "port", "8080")
		isReady.Store(true)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	<-shutdown
	isReady.Store(false)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next(rw, r)
		httpRequestsTotal.WithLabelValues(r.URL.Path, fmt.Sprintf("%d", rw.statusCode)).Inc()
	}
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}
func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

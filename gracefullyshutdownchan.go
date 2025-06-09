package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests received",
}, []string{"status", "path", "method"})

// Middleware to count HTTP requests
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the ResponseWriter to capture the status code
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(recorder, r)

		method := r.Method
		path := r.URL.Path // Path can be adjusted for aggregation (e.g., `/users/:id` → `/users/{id}`)
		status := strconv.Itoa(recorder.statusCode)

		// Increment the counter
		httpRequestCounter.WithLabelValues(status, path, method).Inc()
	})
}

func main() {
	mux := http.NewServeMux()

	reg := prometheus.NewRegistry()

	reg.MustRegister(httpRequestCounter)

	handler := promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{})

	mux.Handle("/metrics", handler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	promHandler := prometheusMiddleware(mux)

	log.Println("Starting HTTP server on port", port)

	if err := http.ListenAndServe(":"+port, promHandler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// Helper to capture HTTP status codes
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

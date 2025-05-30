package prometheus

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestedTypes = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requested_types_total",
		Help: "Number of requests by type",
	}, []string{"type"})

	RequestedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "requested_total",
		Help: "Total number of requests",
	})

	ResponseTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "response_time",
		Help: "Response time in nanoseconds",
	}, []string{"path"})

	ImageRequest = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "image_request",
		Help: "Number of requests with image requested or not",
	}, []string{"status"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration",
		Help:    "Duration of request processing in nanoseconds",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1.0, 2.5, 5.0, 7.5, 10.0},
	}, []string{"trafficlight", "need_image"})

	ErrorsAmount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "errors_amount_total",
		Help: "Http errors",
	}, []string{"error"})
)

func init() {
	prometheus.MustRegister(ErrorsAmount)
}

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Nanoseconds()
		ResponseTime.WithLabelValues(r.URL.Path).Set(float64(duration))
	})
}

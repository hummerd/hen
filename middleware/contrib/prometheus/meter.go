package henprometheus

import (
	"net/http"
	"strconv"
	"time"

	"github.com/hummerd/hen"
	"github.com/prometheus/client_golang/prometheus"
)

func DefaultHTTPRequestHistogram() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30, 60},
	}, []string{"method", "endpoint", "status"})
}

func Meter(hist *prometheus.HistogramVec) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			var err error

			ctx := req.Context()

			endpoint := hen.EndpointFromContext(ctx)
			if endpoint == "" {
				return next(req)
			}

			s := time.Now()

			resp, err := next(req)

			status := "error"
			if resp != nil {
				status = strconv.Itoa(resp.StatusCode)
			}

			hist.
				WithLabelValues(
					req.Method,
					endpoint,
					status).
				Observe(time.Since(s).Seconds())

			if err != nil {
				return resp, err
			}

			return resp, nil
		}
	}
}

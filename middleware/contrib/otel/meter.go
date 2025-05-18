package henotel

import (
	"net/http"
	"strconv"
	"time"

	"github.com/hummerd/hen"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func Meter(trace trace.Tracer, meter metric.Meter) hen.Middleware {
	latencyMeasure, err := meter.Float64Histogram(
		"http.client.duration",
		metric.WithUnit("ms"),
		metric.WithDescription("Measures the duration of outbound HTTP requests."),
	)
	if err != nil {
		panic(err)
	}

	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			var err error

			ctx := req.Context()

			endpoint := hen.EndpointFromContext(ctx)

			s := time.Now()

			resp, err := next(req)

			durationMs := float64(time.Since(s)) / float64(time.Millisecond)

			status := "error"
			if resp != nil {
				status = strconv.Itoa(resp.StatusCode)
			}

			latencyMeasure.Record(ctx, durationMs,
				metric.WithAttributes(
					attribute.String("http.client.method", req.Method),
					attribute.String("http.client.endpoint", endpoint),
					attribute.String("http.client.status", status),
				))

			if err != nil {
				return resp, err
			}

			return resp, nil
		}
	}
}

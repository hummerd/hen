package henotel

import (
	"net/http"

	"github.com/hummerd/hen"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

func Tracer(tracer trace.Tracer) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			var err error

			ctx := req.Context()

			endpoint := hen.EndpointFromContext(ctx)

			ctx, span := tracer.Start(ctx, "http.client.request",
				trace.WithAttributes(
					semconv.HTTPMethod(req.Method),
					semconv.HTTPURL(req.URL.String()),
					attribute.String("http.client.endpoint", endpoint),
				))

			defer func() {
				if err != nil {
					span.RecordError(err)
				}

				// span.SetStatus()
				span.End()
			}()

			resp, err := next(req)
			if err != nil {
				return resp, err
			}

			span.SetAttributes(
				semconv.HTTPStatusCode(resp.StatusCode),
			)

			return resp, nil
		}
	}
}

package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/hummerd/hen"
)

func SLogger(log *slog.Logger, normalLevel, warnLevel, errLevel slog.Level) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()

			log.Log(ctx, normalLevel, "http request",
				"method", req.Method,
				"url", req.URL,
			)

			endpoint := hen.EndpointFromContext(ctx)

			start := time.Now()

			resp, err := next(req)
			if err != nil {
				log.Log(ctx, errLevel, "http response",
					"method", req.Method,
					"url", req.URL,
					"endpoint", endpoint,
					"error", err,
				)

				return resp, err
			}

			args := []any{
				"method", req.Method,
				"url", req.URL,
				"endpoint", endpoint,
				"status", resp.Status,
				"duration", time.Since(start),
			}

			switch {
			case resp.StatusCode >= http.StatusInternalServerError:
				log.Log(ctx, errLevel, "http response", args...)
			case resp.StatusCode >= http.StatusBadRequest:
				log.Log(ctx, warnLevel, "http response", args...)
			default:
				log.Log(ctx, normalLevel, "http response", args...)
			}

			return resp, nil
		}
	}
}

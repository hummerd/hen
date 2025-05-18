package zap

import (
	"net/http"
	"time"

	"github.com/hummerd/hen"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger(log *zap.Logger, normalLevel, warnLevel, errorLevel zapcore.Level) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			log.Log(normalLevel, "http request",
				zap.String("method", req.Method),
				zap.Stringer("url", req.URL))

			start := time.Now()

			resp, err := next(req)
			if err != nil {
				log.Log(errorLevel, "http response",
					zap.String("method", req.Method),
					zap.Stringer("url", req.URL),
					zap.String("status", resp.Status),
					zap.Duration("duration", time.Since(start)),
					zap.Error(err))

				return resp, err
			}

			var level zapcore.Level

			switch {
			case resp.StatusCode >= http.StatusInternalServerError:
				level = errorLevel
			case resp.StatusCode >= http.StatusBadRequest:
				level = warnLevel
			default:
				level = normalLevel
			}

			log.Log(level, "http response",
				zap.String("method", req.Method),
				zap.Stringer("url", req.URL),
				zap.String("status", resp.Status),
				zap.Duration("duration", time.Since(start)))

			return resp, nil
		}
	}
}

package zerolog

import (
	"net/http"
	"time"

	"github.com/hummerd/hen"
	"github.com/rs/zerolog"
)

func Logger(log *zerolog.Logger, normalLevel, warnLevel, errorLevel zerolog.Level) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			log.WithLevel(normalLevel).
				Str("method", req.Method).
				Stringer("url", req.URL).
				Msg("http request")

			start := time.Now()

			resp, err := next(req)
			if err != nil {
				log.WithLevel(errorLevel).
					Str("method", req.Method).
					Stringer("url", req.URL).
					Err(err).
					Msg("http response")

				return resp, err
			}

			var le *zerolog.Event

			switch {
			case resp.StatusCode >= http.StatusInternalServerError:
				le = log.WithLevel(errorLevel)
			case resp.StatusCode >= http.StatusBadRequest:
				le = log.WithLevel(warnLevel)
			default:
				le = log.WithLevel(normalLevel)
			}

			le.Str("method", req.Method).
				Stringer("url", req.URL).
				Str("status", resp.Status).
				Dur("duration", time.Since(start)).
				Msg("http response")

			return resp, nil
		}
	}
}

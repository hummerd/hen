package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/hummerd/hen"
)

func Logger(log *log.Logger) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			log.Printf("http request: %s %s", req.Method, req.URL)

			start := time.Now()

			resp, err := next(req)
			if err != nil {
				log.Printf("http response: %s %s %v", req.Method, req.URL, err)
				return resp, err
			}

			log.Printf("http response: %s %s %s %v", req.Method, req.URL, resp.Status, time.Since(start))
			return resp, nil
		}
	}
}

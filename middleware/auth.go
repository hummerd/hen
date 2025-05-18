package middleware

import (
	"net/http"

	"github.com/hummerd/hen"
)

func BasicAuth(username, password string) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			req.SetBasicAuth(username, password)
			return next(req)
		}
	}
}

func Auth(f func(req *http.Request) string) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			req.Header.Set(hen.HeaderAuthorization, f(req))
			return next(req)
		}
	}
}

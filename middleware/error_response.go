package middleware

import (
	"io"
	"net/http"

	"github.com/hummerd/hen"
)

func HTTPErrorResponse(successCodes ...int) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			resp, err := next(req)
			if err != nil {
				return resp, err
			}

			for _, code := range successCodes {
				if resp.StatusCode == code {
					return resp, nil
				}
			}

			defer resp.Body.Close()

			endpoint := hen.EndpointFromContext(req.Context())

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return resp, err
			}

			return resp, &hen.HTTPError{
				URL:        req.URL,
				Endpoint:   endpoint,
				Method:     req.Method,
				StatusCode: resp.StatusCode,
				Body:       respBody,
			}
		}
	}
}

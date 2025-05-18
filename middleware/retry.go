package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/hummerd/hen"
)

type RetryContext struct {
	Request       *http.Request
	Response      *http.Response
	Error         error
	Attempt       int
	TotalDuration time.Duration
	LastDuration  time.Duration
}

type RetryFunc func(retryContext RetryContext) (bool, time.Duration)

func ExponentialBackOff(minDelay, maxDelay time.Duration) RetryFunc {
	return func(retryContext RetryContext) (bool, time.Duration) {
		needRetry := retryContext.Error != nil ||
			retryContext.Response.StatusCode >= http.StatusInternalServerError

		delay := minDelay << uint(retryContext.Attempt-1)
		if delay > maxDelay {
			delay = maxDelay
		}

		return needRetry, delay
	}
}

func ConstBackOff(delay time.Duration) RetryFunc {
	return func(retryContext RetryContext) (bool, time.Duration) {
		needRetry := retryContext.Error != nil ||
			retryContext.Response.StatusCode >= http.StatusInternalServerError
		return needRetry, delay
	}
}

// Retry is a middleware that retries a request a number of times with a delay between each retry.
// The delay between each retry is calculated by the backoff function.
// Note: To be able to perform retry the request body must be rewindable. If it is not rewindable
// Retry will have to memorise the body before first request and then use it for each retry. Be caution
// if you send large files or so, although it should not be problem for regular jsons or small texts.
func Retry(
	maxRetries int,
	retryFunc RetryFunc,
) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			newGetBody, err := getGetBody(req)
			if err != nil {
				return nil, err
			}

			req.GetBody = newGetBody

			ctx := req.Context()

			beginning := time.Now()

			for i := 0; ; i++ {
				start := time.Now()

				resp, err := next(req)

				retryContext := RetryContext{
					Request:       req,
					Response:      resp,
					Error:         err,
					Attempt:       i + 1,
					TotalDuration: time.Since(beginning),
					LastDuration:  time.Since(start),
				}
				needRetry, delay := retryFunc(retryContext)

				if !needRetry {
					return resp, err
				}

				if i > maxRetries {
					return resp, err
				}

				resp.Body.Close()

				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
				}

				body, err := req.GetBody()
				if err != nil {
					return nil, err
				}

				req = req.Clone(ctx)
				req.Body = body
			}
		}
	}
}

func getGetBody(req *http.Request) (func() (io.ReadCloser, error), error) {
	if req.Body == nil {
		return func() (io.ReadCloser, error) {
			return nil, nil
		}, nil
	}

	if req.GetBody != nil {
		return req.GetBody, nil
	}

	var buff bytes.Buffer // add sync pool here

	_, err := buff.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}

	return func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(buff.Bytes())), nil
	}, nil
}

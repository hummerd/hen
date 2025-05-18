package middleware

import (
	"net/http"
	"testing"
	"time"
)

func TestRetryPolicy(t *testing.T) {
	tests := []struct {
		name       string
		retryFunc  RetryFunc
		retryCount int
		wantRetry  bool
	}{
		{
			name:       "ExponentialBackOff with error",
			retryFunc:  ExponentialBackOff(1*time.Second, 10*time.Second),
			retryCount: 3,
			wantRetry:  true,
		},
		{
			name:       "ConstBackOff with error",
			retryFunc:  ConstBackOff(2 * time.Second),
			retryCount: 3,
			wantRetry:  true,
		},
		{
			name:       "ExponentialBackOff with success",
			retryFunc:  ExponentialBackOff(1*time.Second, 10*time.Second),
			retryCount: 0,
			wantRetry:  false,
		},
		{
			name:       "ConstBackOff with success",
			retryFunc:  ConstBackOff(2 * time.Second),
			retryCount: 0,
			wantRetry:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryContext := RetryContext{
				Request:       &http.Request{},
				Response:      &http.Response{StatusCode: http.StatusInternalServerError},
				Error:         nil,
				Attempt:       tt.retryCount,
				TotalDuration: 0,
				LastDuration:  0,
			}

			gotRetry, _ := tt.retryFunc(retryContext)
			if gotRetry != tt.wantRetry {
				t.Errorf("RetryFunc() = %v, want %v", gotRetry, tt.wantRetry)
			}
		})
	}
}

func TestRetryMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		retryFunc  RetryFunc
		maxRetries int
		wantErr    bool
	}{
		{
			name:       "Retry with ExponentialBackOff",
			retryFunc:  ExponentialBackOff(1*time.Second, 10*time.Second),
			maxRetries: 3,
			wantErr:    false,
		},
		{
			name:       "Retry with ConstBackOff",
			retryFunc:  ConstBackOff(2 * time.Second),
			maxRetries: 3,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := Retry(tt.maxRetries, tt.retryFunc)
			next := func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusInternalServerError}, nil
			}

			req, _ := http.NewRequest("GET", "http://example.com", nil)
			_, err := middleware(next)(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

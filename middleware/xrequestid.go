package middleware

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/hummerd/hen"
)

func XRequestID(f func(req *http.Request) string) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			req.Header.Set(hen.HeaderXRequestID, f(req))
			return next(req)
		}
	}
}

func XRequestIDFromContext(key any) hen.Middleware {
	return func(next hen.DoFunc) hen.DoFunc {
		return func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()

			reqID := ""
			if v, ok := ctx.Value(key).(string); ok {
				reqID = v
			} else {
				reqID = randString()
			}

			req.Header.Set(hen.HeaderXRequestID, reqID)

			return next(req)
		}
	}
}

func randString() string {
	t := uint64(time.Now().UnixNano())
	gen := rand.NewPCG(t, t/7)

	a := gen.Uint64()
	b := gen.Uint64()

	buff := make([]byte, 0, 16)

	binary.BigEndian.AppendUint64(buff, a)
	binary.BigEndian.AppendUint64(buff, b)

	return base64.URLEncoding.EncodeToString(buff)
}

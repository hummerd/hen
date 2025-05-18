package hen

import (
	"net/http"
)

type DoFunc func(req *http.Request) (*http.Response, error)

type Middleware func(next DoFunc) DoFunc

package hen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
)

func NewHen(client *http.Client, baseURL string, middlewares ...Middleware) *Hen {
	if baseURL != "" && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	return &Hen{
		client:      client,
		baseURL:     baseURL,
		middlewares: middlewares,
	}
}

type Hen struct {
	client      *http.Client
	baseURL     string
	middlewares []Middleware
}

func (h *Hen) NewEndpoint(path, method string, middlewares ...Middleware) *Endpoint {
	if path != "" && path[0] != '/' {
		path = "/" + path
	}

	e := &Endpoint{
		hen:    h,
		path:   path,
		method: method,
	}

	e.doFunc = e.do

	for _, m := range h.middlewares {
		e.doFunc = m(e.doFunc)
	}

	for _, m := range middlewares {
		e.doFunc = m(e.doFunc)
	}

	return e
}

type Endpoint struct {
	hen    *Hen
	path   string
	method string
	// doFunc is a function that performs the actual HTTP request. It is Endpoint.do() function
	// wrapped by middlewares.
	doFunc DoFunc
}

type henContextKey int

const (
	ContextKeyEndpoint henContextKey = iota
)

func EndpointFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ContextKeyEndpoint).(string)
	return v
}

func (e *Endpoint) DoJSON(
	ctx context.Context,
	pathArgs []any,
	query urlpkg.Values,
	requestObject any,
	responseObject any,
) error {
	body, header, err := jsonBody(requestObject)
	if err != nil {
		return err // add more details here
	}

	resp, err := e.Do(ctx, pathArgs, query, header, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if responseObject != nil {
		err = json.NewDecoder(resp.Body).Decode(responseObject)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err) // add more details here
		}
	}

	return nil
}

func jsonBody(requestObject any) (io.Reader, http.Header, error) {
	if requestObject == nil {
		return nil, nil, nil
	}

	var buff bytes.Buffer // add sync pool here

	header := http.Header{}
	header.Set(HeaderContentType, "application/json; charset=utf-8")

	err := json.NewEncoder(&buff).Encode(requestObject)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request body: %w", err) // add more details here
	}

	return &buff, header, nil
}

func (e *Endpoint) Do(
	ctx context.Context,
	pathArgs []any,
	query urlpkg.Values,
	headers http.Header,
	body io.Reader,
) (*http.Response, error) {
	url := e.preparePath(pathArgs, query)

	ctx = context.WithValue(ctx, ContextKeyEndpoint, e.path)

	req, err := http.NewRequestWithContext(ctx, e.method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range headers {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	return e.doFunc(req)
}

func (e *Endpoint) do(r *http.Request) (*http.Response, error) {
	return e.hen.client.Do(r)
}

func (e *Endpoint) preparePath(pathArgs []any, query urlpkg.Values) string {
	escapeArgs := make([]any, len(pathArgs))

	for i, arg := range pathArgs {
		if str, ok := arg.(string); ok {
			escapeArgs[i] = urlpkg.PathEscape(str)
		} else {
			escapeArgs[i] = arg
		}
	}

	url := e.hen.baseURL + e.path
	if len(pathArgs) > 0 { // TODO: what if there is not enough args?
		url = fmt.Sprintf(url, escapeArgs...)
	}

	if len(query) > 0 {
		url = url + "?" + query.Encode()
	}

	return url
}

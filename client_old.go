package hen

import (
	"fmt"
	"io"
	"net/url"
)

// func NewClient(client *http.Client, baseURL string) *HClient {
// 	return &HClient{
// 		c:       client,
// 		baseURL: baseURL,
// 	}
// }

// type HClient struct {
// 	c       *http.Client
// 	baseURL string
// }

// type HeaderFunc func() string

// func (c *HClient) NewEndpoint(path, method string, auth HeaderFunc) *Endpoint {
// 	return &Endpoint{
// 		c:      c,
// 		path:   path,
// 		method: method,
// 		auth:   auth,
// 	}
// }

// type Endpoint struct {
// 	c      *HClient
// 	path   string
// 	method string
// 	// headers []httpHeader
// 	encoder func(interface{}) (io.ReadCloser, string, error)
// 	decoder func(io.Reader, interface{}) error
// 	auth    HeaderFunc
// 	mw      Middleware
// }

// func (e *Endpoint) DoJSON(
// 	ctx context.Context,
// 	pathArgs []interface{},
// 	query url.Values,
// 	request, response interface{},
// ) error {
// 	var reqBody io.Reader

// 	var headers []httpHeader

// 	if request != nil {
// 		var buff bytes.Buffer

// 		err := json.NewEncoder(&buff).Encode(request)
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal request body: %w", err)
// 		}

// 		reqBody = &buff
// 		headers = append(headers, httpHeader{
// 			name:  "Content-Type",
// 			value: "application/json; charset=utf-8",
// 		})
// 	}

// 	resp, err := e.do(ctx, pathArgs, query, headers, reqBody)
// 	if err != nil {
// 		return fmt.Errorf("failed to perform request: %w", err)
// 	}

// 	defer resp.Body.Close()

// 	if response != nil {
// 		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
// 			return fmt.Errorf("failed to decode response body: %w", err)
// 		}
// 	}

// 	return nil
// }

type HTTPError struct {
	URL        *url.URL
	Endpoint   string
	Method     string
	StatusCode int
	Body       []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf(
		"request failed: %s %s: %d",
		e.Method,
		e.Endpoint,
		e.StatusCode,
	)
}

type Args struct {
	PathArgs       []interface{}
	Query          url.Values
	RequestObject  interface{}
	RawRequestBody io.Reader
	ResponseObject interface{}
}

// func (e *Endpoint) Do(ctx context.Context, args *Args) (*http.Response, error) {
// 	return e.do(ctx,

// 		pathArgs, query, nil, body)
// }

// func (e *Endpoint) Do(ctx context.Context, args *Args) (*http.Response, error) {
// 	return e.do(ctx, args)
// }

// type httpHeader struct {
// 	name  string
// 	value string
// }

// func (e *Endpoint) do(
// 	ctx context.Context,
// 	args *Args,
// ) (*http.Response, error) {
// 	url := e.preparePath(args.PathArgs, args.Query)

// 	req, err := http.NewRequestWithContext(ctx, e.method, url, args.RawRequestBody)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}

// 	for _, h := range e.headers {
// 		req.Header.Set(h.name, h.value)
// 	}

// 	if e.encoder != nil && args.RequestObject != nil {
// 		body, contentType, err := e.encoder(args.RequestObject)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to encode body: %w", err)
// 		}

// 		req.Body = body
// 		req.Header.Set("Content-Type", contentType)
// 	}

// 	if e.auth != nil {
// 		req.Header.Set("Authorization", e.auth())
// 	}

// 	var resp *http.Response

// 	// do loop to allow scenarios like:
// 	// retry on unauthorized
// 	// retry on error
// 	retry := true
// 	retryCount := 0
// 	for retry {
// 		err = e.mw.Before(ctx, PreRequestData{
// 			Path:       e.path,
// 			Request:    req,
// 			RetryCount: retryCount,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}

// 		s := time.Now()

// 		resp, err = e.c.c.Do(req)

// 		if err == nil && e.decoder != nil && args.RequestObject != nil {
// 			err = e.decoder(resp.Body, args.ResponseObject)
// 			if err != nil {
// 				err = fmt.Errorf("failed to decode response body: %w", err)
// 			}
// 		}

// 		retry, err = e.mw.After(ctx, PostRequestData{
// 			Path:       e.path,
// 			Duration:   time.Since(s),
// 			Request:    req,
// 			Response:   resp,
// 			RetryCount: retryCount,
// 		}, err)
// 		if err != nil {
// 			if resp != nil {
// 				resp.Body.Close()
// 			}
// 			return nil, err
// 		}

// 		retryCount++
// 	}

// 	return resp, err
// }

// func (e *Endpoint) preparePath(pathArgs []interface{}, query url.Values) string {
// 	url := e.path
// 	if len(pathArgs) > 0 {
// 		url = path.Join(e.c.baseURL, e.path)
// 		url = fmt.Sprintf(url, pathArgs...)
// 	}

// 	if len(query) > 0 {
// 		url = url + "?" + query.Encode()
// 	}

// 	return url
// }

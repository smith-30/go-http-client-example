package restclient

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Request struct {
	c *RESTClient

	timeout time.Duration

	verb       string
	pathPrefix string
	subpath    string
	params     url.Values
	headers    http.Header

	err  error
	body io.Reader
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestConstructionError struct {
	Err error
}

func (r *RequestConstructionError) Error() string {
	return fmt.Sprintf("request construction error: '%v'", r.Err)
}

func NewRequest(c *RESTClient) *Request {
	var pathPrefix string
	if c.base != nil {
		pathPrefix = path.Join("/", c.base.Path, c.versionedAPIPath)
	} else {
		pathPrefix = path.Join("/", c.versionedAPIPath)
	}

	var timeout time.Duration
	if c.Client != nil {
		timeout = c.Client.Timeout
	}

	r := &Request{
		c:          c,
		backoff:    backoff,
		timeout:    timeout,
		pathPrefix: pathPrefix,
	}

	switch {
	case len(c.content.AcceptContentTypes) > 0:
		r.SetHeader("Accept", c.content.AcceptContentTypes)
	case len(c.content.ContentType) > 0:
		r.SetHeader("Accept", c.content.ContentType+", */*")
	}
	return r
}

func (r *Request) SetHeader(key string, values ...string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Del(key)
	for _, value := range values {
		r.headers.Add(key, value)
	}
	return r
}

func (r *Request) Timeout(d time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = d
	return r
}

// DoRaw executes the request but does not process the response body.
func (r *Request) DoRaw(ctx context.Context) ([]byte, error) {
	var result Result
	err := r.request(ctx, func(req *http.Request, resp *http.Response) {
		result.body, result.err = ioutil.ReadAll(resp.Body)
		if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
			// Todo error body のパース
			result.err = &StatusError{Message: fmt.Sprintf("%v", resp.StatusCode)}
		}
	})
	if err != nil {
		return nil, err
	}
	return result.body, result.err
}

func (r *Request) request(ctx context.Context, fn func(*http.Request, *http.Response)) error {
	if r.err != nil {
		return r.err
	}

	client := r.c.Client
	if client == nil {
		client = http.DefaultClient
	}

	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	// Right now we make about ten retry attempts if we get a Retry-After response.
	req, err := r.newHTTPRequest(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	defer readAndCloseResponseBody(resp)

	// if the the server returns an error in err, the response will be nil.
	f := func(req *http.Request, resp *http.Response) {
		if resp == nil {
			return
		}
		fn(req, resp)
	}

	f(req, resp)

	return err
}

func (r *Request) newHTTPRequest(ctx context.Context) (*http.Request, error) {
	url := r.URL().String()
	req, err := http.NewRequest(r.verb, url, r.body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header = r.headers
	return req, nil
}

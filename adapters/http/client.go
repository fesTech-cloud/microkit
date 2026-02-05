package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/festech-cloud/microkit/network"
	"github.com/festech-cloud/microkit/internal/retry"
)

type Client struct {
	client *http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Get(ctx context.Context, url string, opts ...network.Option) (*network.Response, error) {
	return c.do(ctx, "GET", url, nil, opts...)
}

func (c *Client) Post(ctx context.Context, url string, body []byte, opts ...network.Option) (*network.Response, error) {
	return c.do(ctx, "POST", url, body, opts...)
}

func (c *Client) Put(ctx context.Context, url string, body []byte, opts ...network.Option) (*network.Response, error) {
	return c.do(ctx, "PUT", url, body, opts...)
}

func (c *Client) Patch(ctx context.Context, url string, body []byte, opts ...network.Option) (*network.Response, error) {
	return c.do(ctx, "PATCH", url, body, opts...)
}

func (c *Client) Delete(ctx context.Context, url string, opts ...network.Option) (*network.Response, error) {
	return c.do(ctx, "DELETE", url, nil, opts...)
}

func (c *Client) Call(ctx context.Context, method string, req interface{}, resp interface{}, opts ...network.Option) error {
	panic("gRPC not supported by HTTP client")
}

func (c *Client) do(ctx context.Context, method, url string, body []byte, opts ...network.Option) (*network.Response, error) {
	config := &network.Config{}
	for _, opt := range opts {
		opt(config)
	}

	if config.RetryConfig != nil {
		var resp *network.Response
		err := retry.Execute(ctx, *config.RetryConfig, func() error {
			var err error
			resp, err = c.doRequest(ctx, method, url, body, config)
			return err
		})
		return resp, err
	}

	return c.doRequest(ctx, method, url, body, config)
}

func (c *Client) doRequest(ctx context.Context, method, url string, body []byte, config *network.Config) (*network.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &network.Response{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       respBody,
	}, nil
}

func (c *Client) Close() error {
	c.client.CloseIdleConnections()
	return nil
}

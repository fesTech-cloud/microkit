package grpc

import (
	"context"
	"time"

	"github.com/festech-cloud/microkit/network"
	"github.com/festech-cloud/microkit/internal/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(target string, timeout time.Duration) (*Client, error) {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Get(ctx context.Context, url string, opts ...network.Option) (*network.Response, error) {
	panic("HTTP methods not supported by gRPC client")
}

func (c *Client) Post(ctx context.Context, url string, body []byte, opts ...network.Option) (*network.Response, error) {
	panic("HTTP methods not supported by gRPC client")
}

func (c *Client) Put(ctx context.Context, url string, body []byte, opts ...network.Option) (*network.Response, error) {
	panic("HTTP methods not supported by gRPC client")
}

func (c *Client) Delete(ctx context.Context, url string, opts ...network.Option) (*network.Response, error) {
	panic("HTTP methods not supported by gRPC client")
}

func (c *Client) Patch(ctx context.Context, url string, body []byte, opts ...network.Option) (*network.Response, error) {
	panic("HTTP methods not supported by gRPC client")
}

func (c *Client) Call(ctx context.Context, method string, req any, resp any, opts ...network.Option) error {
	config := &network.Config{}
	for _, opt := range opts {
		opt(config)
	}

	if config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}

	if config.RetryConfig != nil {
		return retry.Execute(ctx, *config.RetryConfig, func() error {
			return c.conn.Invoke(ctx, method, req, resp)
		})
	}

	return c.conn.Invoke(ctx, method, req, resp)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

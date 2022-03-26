package httpclient

import "time"

type Option func(*Client)

func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithRetrier(retrier string) Option {
	return func(c *Client) {
		c.retrier = retrier
	}
}

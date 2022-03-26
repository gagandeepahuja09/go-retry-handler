package httpclient

import (
	"net/http"
	"time"

	"github.com/gagandeepahuja09/goretryhandler"
)

const (
	defaultHTTPTimeout = 30 * time.Second
	defaultRetryCount  = 0
)

type Client struct {
	client goretryhandler.Doer

	timeout    time.Duration
	retryCount int
	retrier    goretryhandler.Retriable
	plugins    []goretryhandler.Plugin
}

func NewClient(opts ...Option) *Client {
	client := Client{
		timeout:    defaultHTTPTimeout,
		retryCount: defaultRetryCount,
		retrier:    goretryhandler.NewNoRetrier(),
	}

	for _, opt := range opts {
		opt(&client)
	}

	if client.client == nil {
		client.client = &http.Client{
			Timeout: client.timeout,
		}
	}

	return &client
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {

}

func (c *Client) reportRequestStart(request *http.Request) {
	for _, plugin := range c.plugins {
		plugin.OnRequestStart(request)
	}
}

func (c *Client) reportRequestEnd(request *http.Request) {
	for _, plugin := range c.plugins {
		plugin.OnRequestEnd(request)
	}
}

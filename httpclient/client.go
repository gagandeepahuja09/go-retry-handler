package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gagandeepahuja09/goretryhandler"
	"github.com/gojek/valkyrie"
	"github.com/gojektech/valkyrie"
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

// Do makes a HTTP request with the native `http.Do` interface
// Includes retries with backoffs.
// Return a multi error if multiple requests fail.
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	request.Close = true

	var bodyReader *bytes.Reader

	if request.Body != nil {
		reqData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqData)
		request.Body = ioutil.NopCloser(bodyReader) // prevents closing the body between retries
	}

	multiErr := &valkyrie.MultiError{}
	var response *http.Response

	for i := 0; i <= c.retryCount; i++ {
		if response != nil {
			response.Body.Close()
		}

		c.reportRequestStart(request)
		var err error
		response, err = c.client.Do(request)
		if bodyReader != nil {
			// Reset the body after the request, since at this point it is already read.
			bodyReader.Seek(0, 0)
		}

		if err != nil {
			multiErr.Push(err.Error())
			c.reportError(request, err)
			backOffTime := c.retrier.NextInterval()
			time.Sleep(backOffTime)
			continue
		}
		c.reportRequestEnd(request, response)

		if response.StatusCode >= http.StatusInternalServerError {
			backOffTime := c.retrier.NextInterval(i)
			time.Sleep(backOffTime)
			continue
		}

		multiErr = &valkyrie.MultiError{} // Clears err if any iteration succeeds
		break
	}

	return response, multiErr.HasError()
}

func (c *Client) AddPlugin(p goretryhandler.Plugin) {
	c.plugins = append(c.plugins, p)
}

func (c *Client) reportRequestStart(request *http.Request) {
	for _, plugin := range c.plugins {
		plugin.OnRequestStart(request)
	}
}

func (c *Client) reportError(request *http.Request, err error) {
	for _, plugin := range c.plugins {
		plugin.OnError(request, err)
	}
}

func (c *Client) reportRequestEnd(request *http.Request, response *http.Response) {
	for _, plugin := range c.plugins {
		plugin.OnRequestEnd(request, response)
	}
}

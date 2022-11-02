package httpmock

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

type CallOption func(*request)

func ReturnError(err error) CallOption {
	return func(r *request) {
		r.returnError = err
	}
}

func ReturnStatus(status int) CallOption {
	return func(r *request) {
		r.returnStatus = status
	}
}

func ReturnBodyRaw(body string) CallOption {
	return func(r *request) {
		r.returnBody = body
	}
}

func ReturnBodyFromObject(object interface{}) CallOption {
	return func(r *request) {
		body, _ := json.Marshal(&object)
		r.returnBody = string(body)
	}
}

func ExpectBody(expectedBody string) CallOption {
	return func(r *request) {
		r.expectedBody = expectedBody
	}
}

func ExpectJSON(expectedJSON string) CallOption {
	return func(r *request) {
		r.expectedJSON = []byte(expectedJSON)
	}
}

func ExpectHeader(name string, values []string) CallOption {
	return func(c *request) {
		if c.expectedHeaders == nil {
			c.expectedHeaders = make(map[string][]string)
		}
		c.expectedHeaders[name] = values
	}
}

func ExpectQueryParamValues(name string, values []string) CallOption {
	return func(c *request) {
		if c.expectedQueryParams == nil {
			c.expectedQueryParams = make(url.Values)
		}
		c.expectedQueryParams[name] = values
	}
}

func ExpectQueryParam(name, value string) CallOption {
	return func(c *request) {
		if c.expectedQueryParams == nil {
			c.expectedQueryParams = make(url.Values)
		}
		c.expectedQueryParams[name] = []string{value}
	}
}

type Client struct {
	http.Client
	transport *transport
}

func (c *Client) WithRequest(method, route string, options ...CallOption) *Client {
	req := newMockRequest(method, route)
	for _, option := range options {
		option(&req)
	}
	c.transport.requests = append(c.transport.requests, req)
	return c
}

func (c *Client) AssertExpectations(t *testing.T) {
	if c.transport.index < len(c.transport.requests) {
		t.Errorf("httpmock should have more requests")
	}
}

func New(t *testing.T) *Client {
	mockTransport := &transport{
		t:        t,
		requests: make([]request, 0),
	}

	return &Client{
		transport: mockTransport,
		Client: http.Client{
			Transport: mockTransport,
		},
	}
}

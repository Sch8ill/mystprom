// Package client provides a simple http client wrapper.
package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HttpClient struct {
	url     string
	client  http.Client
	headers map[string]string
}

func New(url string) *HttpClient {
	return &HttpClient{
		url:     url,
		headers: map[string]string{},
	}
}

func (c *HttpClient) Get(path string) (*http.Response, error) {
	res, err := c.doRequest(path, "GET", []byte{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *HttpClient) Post(path string, body []byte) (*http.Response, error) {
	return c.doRequest(path, "POST", body)
}

func (c *HttpClient) PostJSON(url string, body any) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return c.Post(url, jsonBody)
}

func (c *HttpClient) SetHeader(key string, value string) {
	c.headers[key] = value
}

func (c *HttpClient) doRequest(path string, method string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, c.fullURL(path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *HttpClient) fullURL(path string) string {
	return c.url + path
}

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/akutz/goof"
)

func (c *Client) httpDo(
	method, path string,
	payload, reply interface{}) (*http.Response, error) {

	reqBody, err := encPayload(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s%s", c.Host, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		req.Header[k] = v
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	c.logResponse(res)

	if res.StatusCode > 299 {
		httpErr, err := goof.DecodeHTTPError(res.Body)
		if err != nil {
			return res, goof.WithField("status", res.StatusCode, "http error")
		}
		return res, httpErr
	}

	if req.Method != http.MethodHead && reply != nil {
		if err := decRes(res.Body, reply); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (c *Client) httpGet(
	path string,
	reply interface{}) (*http.Response, error) {

	return c.httpDo("GET", path, nil, reply)
}

func (c *Client) httpHead(
	path string) (*http.Response, error) {

	return c.httpDo("HEAD", path, nil, nil)
}

func (c *Client) httpPost(
	path string,
	payload interface{},
	reply interface{}) (*http.Response, error) {

	return c.httpDo("POST", path, payload, reply)
}

func (c *Client) httpDelete(
	path string,
	reply interface{}) (*http.Response, error) {

	return c.httpDo("DELETE", path, nil, reply)
}

func encPayload(payload interface{}) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

func decRes(body io.Reader, reply interface{}) error {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buf, reply); err != nil {
		return err
	}
	return nil
}

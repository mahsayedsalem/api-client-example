package form3_api_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	contentType    = "application/vnd.api+json"
)

type client struct {
	baseUrl *url.URL
	httpClient *http.Client
}

type response struct {
	*http.Response
}

type errorResponse struct {
	Response   *http.Response
	StatusCode int
	Code       string `json:"error_code"`
	Message    string `json:"error_message"`
}

func newClient(httpClient *http.Client, baseURLStr string) *client{
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	baseURL, _ := url.Parse(baseURLStr)
	return &client{
		baseUrl: baseURL,
		httpClient: httpClient,
	}
}

func (c *client) get(url string, body interface{}) (*http.Request, error) {
	return c.newRequest(http.MethodGet, url, body)
}

func (c *client) post(url string, body interface{}) (*http.Request, error) {
	return c.newRequest(http.MethodPost, url, body)
}

func (c *client) delete(url string, body interface{}) (*http.Request, error) {
	return c.newRequest(http.MethodDelete, url, body)
}

func (c *client) newRequest(method, urlString string, payload interface{}) (*http.Request, error){
	reqUrl, err := c.baseUrl.Parse(urlString)
	if err != nil {
		return nil, err
	}

	var data io.Reader
	if payload != nil {
		body := struct {
			Data interface{} `json:"data"`
		}{Data: payload}
		dataJson, _ := json.Marshal(body)
		data = bytes.NewBuffer(dataJson)
	}

	req, err := http.NewRequest(method, reqUrl.String(), data)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	return req, nil
}

func (c *client) do(ctx context.Context, req *http.Request, val interface{}) (*response, error) {

	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = checkResponseErrors(resp)
	if err != nil {
		return response, err
	}

	respText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	if val != nil {
		body := &struct {
			Data interface{} `json:"data"`
		}{}
		err = json.Unmarshal(respText, body)
		if err != nil {
			return response, err
		}

		encoded, err := json.Marshal(body.Data)
		if err == nil {
			err = json.Unmarshal(encoded, val)
		}
	}
	return response, err
}

func newResponse(r *http.Response) *response {
	return &response{Response: r}
}

func (e *errorResponse) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.StatusCode, e.Message)
}

func checkResponseErrors(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &errorResponse{Response: r, StatusCode: r.StatusCode}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}
	return errorResponse
}

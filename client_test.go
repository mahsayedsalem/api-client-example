package form3_api_client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupClient() (*mux.Router, *httptest.Server, *client) {
	router := mux.NewRouter()
	server := httptest.NewServer(router)
	return router, server, newClient(nil, server.URL)
}

func testRequestMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct
	{
		description string
		input *http.Client
		expected *http.Client
	}{
		{
			description: "Send nil to create new client",
			input: nil,
			expected: &http.Client{
				Timeout: 10 * time.Second,
			},

		},
		{
			description: "Send an http client",
			input: &http.Client{
				Timeout: 15 * time.Second,
			},
			expected: &http.Client{
				Timeout: 15 * time.Second,
			},
		},
		}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := newClient(test.input, BaseUrlEnvVariable)
			assert.Equal(t, fmt.Sprintf("%T", got.httpClient), fmt.Sprintf("%T", test.expected))
			assert.Equal(t, test.expected.Timeout, got.httpClient.Timeout)
		})
	}
}

func TestDo(t *testing.T) {
	router, server, client := setupClient()
	defer server.Close()

	type account struct {
		Name string
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"Name":"form3"}}`))
	}).Methods(http.MethodGet)

	req, _ := client.newRequest(http.MethodGet, "/", nil)
	body := account{}
	_, err := client.do(context.TODO(), req, &body)
	require.Nil(t, err)
	assert.Equal(t, account{"form3"}, body)
}

func TestClientGET(t *testing.T) {
	c := newClient(nil, BaseUrlEnvVariable)
	req, _ := c.newRequest(http.MethodGet, "/account", nil)
	testRequestMethod(t, req, http.MethodGet)
}

func TestClientPOST(t *testing.T) {
	c := newClient(nil, BaseUrlEnvVariable)
	req, _ := c.newRequest(http.MethodPost, "/account", nil)
	testRequestMethod(t, req, http.MethodPost)
}

func TestClientDELETE(t *testing.T) {
	c := newClient(nil, BaseUrlEnvVariable)
	req, _ := c.newRequest(http.MethodDelete, "/account", nil)
	testRequestMethod(t, req, http.MethodDelete)
}

func TestNewRequest(t *testing.T) {
	c := newClient(nil, BaseUrlEnvVariable)
	type account struct {
		Name string
	}

	body := account{Name: "form3"}
	req, err := c.newRequest(http.MethodGet, "/account", &body)
	require.Nil(t, err)

	expectedURL := fmt.Sprintf("%saccount", BaseUrlEnvVariable)
	assert.Equal(t, expectedURL, req.URL.String())

	reqBody, err := ioutil.ReadAll(req.Body)
	require.Nil(t, err)
	assert.Equal(t, "{\"data\":{\"Name\":\"form3\"}}", string(reqBody))
}

func TestErrorResponse(t *testing.T) {
	tests := []struct
	{
		description            string
		input      *errorResponse
		expected string
	}{
		{
			description: "check Error() string message",
			input: &errorResponse{
				Response:   nil,
				StatusCode: http.StatusBadRequest,
				Message:    "Bad Request",
			},
			expected: "code: 400, message: Bad Request",
		},
		{
			description: "check error message in case of good request",
			input: &errorResponse{},
			expected: "code: 0, message: ",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			assert.Equal(t, test.expected, test.input.Error())
		})
	}
}


func TestCheckResponseError(t *testing.T) {
	tests := []struct
	{
		description         string
		input      			*http.Response
		expected 			*errorResponse
	}{
		{
			description: "check response errors when bad request",
			input: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: ioutil.NopCloser(bytes.NewBufferString("Bad Request")),
			},
			expected: &errorResponse{
				Response:   nil,
				StatusCode: http.StatusBadRequest,
				Message:    "Bad Request",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			assert.Equal(t, test.expected.Error(), checkResponseErrors(test.input).Error())
		})
	}
}

func TestCheckResponseNoError(t *testing.T) {
	input := &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(bytes.NewBufferString("")),
	}
	require.Nil(t, checkResponseErrors(input))
}


package form3_api_client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
	c := newClient(nil, BaseUrlEnvVariable)

	if got, want := c.baseUrl.String(), BaseUrlEnvVariable; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestDo(t *testing.T) {
	router, server, client := setupClient()
	defer server.Close()

	type foo struct {
		A string
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"A":"a"}}`))
	}).Methods(http.MethodGet)

	req, _ := client.newRequest(http.MethodGet, "/", nil)
	body := foo{}
	_, err := client.do(context.TODO(), req, &body)
	require.Nil(t, err)
	assert.Equal(t, foo{"a"}, body)
}

func TestClientGET(t *testing.T) {
	c := newClient(nil, BaseUrlEnvVariable)
	req, _ := c.newRequest("GET", "/foo", nil)
	testRequestMethod(t, req, "GET")
}

func TestClientPOST(t *testing.T) {
	c := newClient(nil, BaseUrlEnvVariable)
	req, _ := c.newRequest("POST", "/foo", nil)
	testRequestMethod(t, req, "POST")
}

func TestClientDELETE(t *testing.T) {
	c := newClient(nil, BaseUrlEnvVariable)
	req, _ := c.newRequest("DELETE", "/foo", nil)
	testRequestMethod(t, req, "DELETE")
}

func TestNewRequest(t *testing.T) {
	c := newClient(nil, BaseUrlEnvVariable)
	type foo struct {
		A string
	}

	body := foo{A: "B"}
	req, err := c.newRequest(http.MethodGet, "/foo", &body)
	require.Nil(t, err)

	expectedURL := fmt.Sprintf("%sfoo", BaseUrlEnvVariable)
	assert.Equal(t, expectedURL, req.URL.String())

	reqBody, err := ioutil.ReadAll(req.Body)
	require.Nil(t, err)
	assert.Equal(t, "{\"data\":{\"A\":\"B\"}}", string(reqBody))
}

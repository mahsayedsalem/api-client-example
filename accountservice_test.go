package form3_api_client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupAccountService() (*mux.Router, *httptest.Server, *AccountService) {
	router := mux.NewRouter()
	server := httptest.NewServer(router)
	return router, server, NewAccountService(server.URL)
}

func TestNewAccountService(t *testing.T) {
	s := NewAccountService(BaseUrlEnvVariable)
	if got, want := s.client.baseUrl.String(), BaseUrlEnvVariable; got != want {
		t.Errorf("NewService BaseURL is %v, want %v", got, want)
	}
}

func TestAccountsServiceCreate(t *testing.T) {
	router, server, service := setupAccountService()
	defer server.Close()

	router.HandleFunc("/v1/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"data":{"id":"1","organisation_id":"2","type":"accounts"}}`))
	}).Methods(http.MethodPost)

	account := AccountData{
		ID: "1",
		OrganisationID: "2",
		Type: "accounts",
	}

	saved, resp, err := service.Create(context.Background(), &account)
	assert.Nil(t, err)

	assert.Equal(t, account.Type, saved.Type)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestAccountsServiceFetch(t *testing.T) {
	router, server, service := setupAccountService()
	defer server.Close()

	router.HandleFunc("/v1/organisation/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"id":"1","organisation_id":"2","type":"accounts"}}`))
	}).Methods(http.MethodGet)

	saved, resp, err := service.Fetch(context.Background(), "1")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "1", saved.ID)
}

func TestAccountsServiceDelete(t *testing.T) {
	router, server, service := setupAccountService()
	defer server.Close()

	router.HandleFunc("/v1/organisation/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete)
	version := int64(1)
	resp, err := service.Delete(context.Background(), "1", &version)
	assert.Nil(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

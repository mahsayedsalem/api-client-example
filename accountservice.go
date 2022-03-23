package form3_api_client

import (
	"context"
	"fmt"
	"net/http"
)

const (
	CreateURL 	   = "v1/organisation/accounts"
	FetchURL  	   = "v1/organisation/accounts/%s"
	DeleteUrl 	   = "v1/organisation/accounts/%s?version=%d"
)

type AccountService struct {
	client *client
}

func NewAccountService(baseURL string) *AccountService{
	return &AccountService{
		client: newClient(nil, BaseUrlEnvVariable),
	}
}

func (s *AccountService) Create(ctx context.Context, account *AccountData) (*AccountData, *response, error) {
	req, err := s.client.newRequest(http.MethodPost, CreateURL, account)
	if err != nil {
		return nil, nil, err
	}
	acc := &AccountData{}
	resp, err := s.client.do(ctx, req, acc)
	if err != nil {
		return nil, resp, err
	}
	return acc, resp, nil
}

func (s *AccountService) Fetch(ctx context.Context, id string) (*AccountData, *response, error) {
	path := fmt.Sprintf(FetchURL, id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	account := &AccountData{}
	resp, err := s.client.do(ctx, req, account)
	if err != nil {
		return nil, resp, err
	}
	return account, resp, err
}

func (s *AccountService) Delete(ctx context.Context, id string, version *int64) (*response, error) {
	path := fmt.Sprintf(DeleteUrl, id, *version)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.do(ctx, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

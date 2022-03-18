package form3_api_client

import (
	"context"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestE2ECreate(t *testing.T) {
	ctx := context.TODO()
	accountService := NewAccountService(BaseUrlEnvVariable)
	country := "GB"
	id, _ := uuid.NewV4()
	organizationId, _ := uuid.NewV4()
	account, res, err := accountService.Create(ctx, &AccountData{
		ID: id.String(),
		OrganisationID: organizationId.String(),
		Type: "accounts",
		Attributes: &AccountAttributes{
			BaseCurrency: "GBP",
			Country: &country,
			Name: []string{"Mahmoud Salem"},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, account.ID, id.String())
	assert.Equal(t, account.OrganisationID, organizationId.String())
}

func TestE2ECreateDuplicate(t *testing.T) {
	ctx := context.TODO()
	accountService := NewAccountService(BaseUrlEnvVariable)
	country := "GB"
	id, _ := uuid.NewV4()
	organizationId, _ := uuid.NewV4()
	account, res, err := accountService.Create(ctx, &AccountData{
		ID: id.String(),
		OrganisationID: organizationId.String(),
		Type: "accounts",
		Attributes: &AccountAttributes{
			BaseCurrency: "GBP",
			Country: &country,
			Name: []string{"Mahmoud Salem"},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, account.ID, id.String())
	assert.Equal(t, account.OrganisationID, organizationId.String())

	_, res, err = accountService.Create(ctx, &AccountData{
		ID: id.String(),
		OrganisationID: organizationId.String(),
		Type: "accounts",
		Attributes: &AccountAttributes{
			BaseCurrency: "GBP",
			Country: &country,
			Name: []string{"Mahmoud Salem"},
		},
	})
	assert.Equal(t, "code: 409, message: Account cannot be created as it violates a duplicate constraint", err.Error())
}

func TestE2EFetch(t *testing.T) {
	ctx := context.TODO()
	accountService := NewAccountService(BaseUrlEnvVariable)
	country := "GB"
	id, _ := uuid.NewV4()
	organizationId, _ := uuid.NewV4()
	account, res, err := accountService.Create(ctx, &AccountData{
		ID: id.String(),
		OrganisationID: organizationId.String(),
		Type: "accounts",
		Attributes: &AccountAttributes{
			BaseCurrency: "GBP",
			Country: &country,
			Name: []string{"Mahmoud Salem"},
		},
	})
	account, res, err = accountService.Fetch(ctx, id.String())
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, account.ID, id.String())
	assert.Equal(t, account.OrganisationID, organizationId.String())
}

func TestE2EFetchResourceNotFound(t *testing.T) {
	ctx := context.TODO()
	accountService := NewAccountService(BaseUrlEnvVariable)
	id, _ := uuid.NewV4()
	_, _, err := accountService.Fetch(ctx, id.String())
	assert.Equal(t, fmt.Sprintf("code: 404, message: record %s does not exist", id), err.Error())
}

func TestE2EDelete(t *testing.T) {
	ctx := context.TODO()
	accountService := NewAccountService(BaseUrlEnvVariable)
	country := "GB"
	id, _ := uuid.NewV4()
	organizationId, _ := uuid.NewV4()
	account, res, err := accountService.Create(ctx, &AccountData{
		ID: id.String(),
		OrganisationID: organizationId.String(),
		Type: "accounts",
		Attributes: &AccountAttributes{
			BaseCurrency: "GBP",
			Country: &country,
			Name: []string{"Mahmoud Salem"},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, account.ID, id.String())
	assert.Equal(t, account.OrganisationID, organizationId.String())
	version := int64(0)
	res, err = accountService.Delete(ctx, id.String(), &version)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestE2EDeleteResourceNotFound(t *testing.T) {
	ctx := context.TODO()
	accountService := NewAccountService(BaseUrlEnvVariable)
	id, _ := uuid.NewV4()
	version := int64(0)
	_, err := accountService.Delete(ctx, id.String(), &version)
	assert.Equal(t, "code: 404, message: ", err.Error())
}


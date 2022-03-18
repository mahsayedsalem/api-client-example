package main

import (
	"context"
	"fmt"
	client "form3-api-client"
)

const BaseUrlEnvVariable = "http://localhost:8000/"

func main() {
	createAccount()
	fetchAccount()
	deleteAccount()
}

func createAccount(){
	ctx := context.TODO()
	accountService := client.NewAccountService(BaseUrlEnvVariable)
	country := "GB"
	account, res, err := accountService.Create(ctx, &client.AccountData{
		ID: "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Type: "accounts",
		Attributes: &client.AccountAttributes{
			BaseCurrency: "GBP",
			Country: &country,
			Name: []string{"Mahmoud Salem"},
		},
	})
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("Created:")
	fmt.Println(account.Attributes.Name)
	fmt.Println(res.Body)
	fmt.Println("________")
}

func fetchAccount(){
	id := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	ctx := context.TODO()
	accountService := client.NewAccountService(BaseUrlEnvVariable)
	account, res, err := accountService.Fetch(ctx, id)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("Fetched: ")
	fmt.Println(account.Attributes.Name)
	fmt.Println(res.Body)
	fmt.Println("________")
}

func deleteAccount(){
	id := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	ctx := context.TODO()
	accountService := client.NewAccountService(BaseUrlEnvVariable)
	version := int64(0)
	res, err := accountService.Delete(ctx, id, &version)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(res.Response.StatusCode)
}

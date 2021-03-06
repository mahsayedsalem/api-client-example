![Unit Tests](https://github.com/mahsayedsalem/form3-api-client/actions/workflows/push.yml/badge.svg)

<h1 align="center">
  Golang API Client Library Example
</h1>

<h4 align="center">An easy-to-use client library for Fake API.</h4>

<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#usage">Usage</a> •
  <a href="#testing">Testing</a> •
  <a href="#example">Example</a> •
  <a href="#continuous-integration">Continuous Integration</a> •
  <a href="#enhancements">Enhancements</a>
</p>

## Key Features

* Create Account using FakeAPI
* Retrieve Account using FakeAPI
* Delete Account using FakeAPI

## Usage

### Run docker-compose to start the FakeAPI server

```sh
$ docker-compose up --build
```

### Import the client library

```
import client "form3-api-client"
```

### Create an account

```
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
```

### Fetch an account
```
id := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
ctx := context.TODO()
accountService := client.NewAccountService(BaseUrlEnvVariable)
account, res, err := accountService.Fetch(ctx, id)
if err != nil{
	fmt.Println(err)
}
```

### Delete an account
```
id := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
ctx := context.TODO()
accountService := client.NewAccountService(BaseUrlEnvVariable)
version := int64(0)
res, err := accountService.Delete(ctx, id, &version)
if err != nil{
	fmt.Println(err)
}
fmt.Println(res.Response.StatusCode)
```

```bigquery
204
```

## Testing
I've included both unit tests and e2e tests. The current coverage is 85%.

### Run Tests

#### Directly
```sh
$ go test
```

#### Inside `docker-compose`
```sh
$ docker-compose logs accountapi-client
```

`e2e_tests` will only pass when the docker-compose services are up and running. Currently I've added the docker-compose build to the workflow so the test cases passes on the CI pipeline.

## Continuous Integration
I created a Github Action workflow to run docker-compose and the test cases with each push. This will make us alert to any faulty code being pushed.

## Example

Visit example folder to run key features directly.

## Enhancements

* Add Validations
* Add List operation with paging

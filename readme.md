# Final Project for udacity golang course
## Description
An restful api that allow users to create, update, delete and get customers
## Setup Instructions
* Install golang onto your machine [install](https://go.dev/doc/install)
* cp .env.example .env
* replace the values in .env with those of your local mysql database
* run the below in your terminal
```
go run ./migrations/migrations.go
go run main.go
```
* go to [localhost:3000](http://localhost:3000)
## Endpoints
* /customers GET
* /customers/:id GET
* /customers/:id DELETE
* /customers/:id PATCH
* /customers POST
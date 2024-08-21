# Bookstore REST API
An implementation of a simple Bookstore REST API using Golang.

## Project Outline
The implementation should follow the rules below:

- Customers should be able to create an account
    - Email should be unique
- Customer should be able to see the list of books
  - Books never run out
- Customer should be able to make an order
  - Don't need to worry about cancellation
  - Don't need to worry about anything related to the payments
  - An order is composed of the arbiter number of books


## Installation 

Clone this repository
```bash
git clone https://github.com/leogsouza/bookstore-api.git
```

Then build and run the services
```bash
$ docker compose up -d --build
```

## API Usage
After the services started the service will be running at localhost:3000 and the following endpoints are available:
* `POST /api/v1/users`: create a new user
* `POST /api/v1/auth/login`: authenticates a user and generates a JWT
* `GET /api/v1/books`: list the books
* `GET /api/v1/users/me/orders`: returns the orders created by the authenticated user


### Todo

- [] Add tests for all APIs
- [] Make docs
- [] Document API with openAPI
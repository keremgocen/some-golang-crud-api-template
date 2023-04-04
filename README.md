# go-product-crud-api

![workflow](https://github.com/keremgocen/some-golang-crud-api-template/actions/workflows/dev.yaml/badge.svg)

## How to Run it?

`-> go run main.go`

### Keyvalue API

#### GET keyvalue/:name

Given a value name, returns the value from storage if it exists.

`/HTTP/GET /keyvalue/:name`

Example call (using curl):
`curl http://localhost:5000/keyvalue/foo`

#### POST keyvalue

Given a value name, returns the value from storage if it exists.

`/HTTP/POST /keyvalue`

Example call (using curl):
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"foo","value":"bar"}' \
  http://localhost:5000/keyvalue
```

### What this API is for?

`keyvalue/service.go` implements:

- GetKeyValue
- PostKeyValue

  functionality for key/value entries to be stored in a local in-memory storage.

Storage layer is injected as a dependency. Current implementation just saves entries in local memory.

#### Error handling

Error handling is very basic, returning appropriate HTTP status codes and a descriptive message.

#### Business logic

Business logic is implemented in keyvalue service via the handlers.

#### API Input validation

Input validation is done by Gin bindings which comes out of boxes, using the default settings.

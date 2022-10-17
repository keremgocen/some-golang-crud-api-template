# go-product-crud-api
![workflow](https://github.com/keremgocen/some-golang-crud-api-template/actions/workflows/dev.yaml/badge.svg)


### What this API is for?

`products/handler.go` implements:

- Create
- List
- Get

  functionality for the product (defined in `models/product.go`) objects which is stored in memory.

There is an optional, mock currency converter that is activated once the `GetProductRequest` has a
`Currency` field set. The default value is assumed as "GBP".

### Running tests locally


### Assumptions

#### Mock Currency Converter

`mockexchange` package has a mock `CurrencyConverter` interface which provides a `ConvertExchangeRate`
function and it only returns a bunch of rates for converting from "GBP".

The idea behind that interface is a live currency converter API can be injected as a depency to the
`productHandler` later on when the product wants to support other currencies. So it abstracts away
the conversion, ignoring empty values or the default value "GBP".

#### Error handling and logging

Error handling is very basic, returning appropriate HTTP status codes and a descriptive message.
Ideally a logger could be injected to the handler to provide a more flexible, informative logging solution.

#### Business logic

Handler is also minimalistic. There is room for more abstraction and moving the business logic outside the
HTTP request handlers.

#### API Schema validation

While I added the JSON encoding tags, schema validation is ignored. It would be a nice next step to
add struct validation for the request handlers.

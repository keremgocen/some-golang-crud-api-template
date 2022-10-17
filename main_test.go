package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/keregocen/go-product-crud-api/pkg/mockexchange"
	"github.com/keregocen/go-product-crud-api/pkg/models"
	"github.com/keregocen/go-product-crud-api/pkg/products"
	"github.com/keregocen/go-product-crud-api/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func run() *mux.Router {
	storageAPI := storage.NewStore()
	currencyConverterAPI := mockexchange.NewConverter()
	service := products.NewService(storageAPI, currencyConverterAPI)

	r := mux.NewRouter()
	r.Handle("/products", middleware(http.HandlerFunc(service.Create))).Methods(http.MethodPost)
	r.Handle("/products", middleware(http.HandlerFunc(service.List))).Methods(http.MethodGet)
	r.Handle("/product", middleware(http.HandlerFunc(service.Get))).Methods(http.MethodPost)
	return r
}

func TestCreateProduct(t *testing.T) {
	router := run()

	req, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name": "New Product"}`))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateProductAlreadyExists(t *testing.T) {
	router := run()

	req, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name": "New Product"}`))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListProduct(t *testing.T) {
	router := run()

	req1, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name": "Product 1"}`))
	assert.Nil(t, err)

	req2, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name": "Product 2"}`))
	assert.Nil(t, err)

	req3, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name": "Product 3"}`))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req1)
	router.ServeHTTP(rr, req2)
	router.ServeHTTP(rr, req3)

	listReq, err := http.NewRequest(http.MethodGet, "/products", nil)
	assert.Nil(t, err)

	listRecorder := httptest.NewRecorder()
	router.ServeHTTP(listRecorder, listReq)
	assert.Equal(t, http.StatusOK, listRecorder.Code)

	var products map[string]models.Product
	err = json.Unmarshal(listRecorder.Body.Bytes(), &products)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(products))
}

func TestGetProduct(t *testing.T) {
	router := run()

	req, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name": "Test product"}`))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	getReq, err := http.NewRequest(http.MethodPost, "/product", strings.NewReader(`{"name": "Test product"}`))
	assert.Nil(t, err)

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getReq)
	assert.Equal(t, http.StatusOK, getRecorder.Code)

	var product models.Product
	err = json.Unmarshal(getRecorder.Body.Bytes(), &product)
	assert.Nil(t, err)

	assert.Equal(t, "Test product", product.Name)
}

func TestGetProductMissing(t *testing.T) {
	router := run()

	req, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name": "Test product"}`))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	getReq, err := http.NewRequest(http.MethodPost, "/product", strings.NewReader(`{"name": "another product"}`))
	assert.Nil(t, err)

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getReq)
	assert.Equal(t, http.StatusNotFound, getRecorder.Code)
}

func TestGetProductWithCurrency(t *testing.T) {
	router := run()

	initialPrice := 100.0
	fromCurrency := "GBP"
	toCurrency := "USD"

	req, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{
		"name": "Test product",
		"price": 100.0,
		"currency": "GBP"
	}`))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	getReq, err := http.NewRequest(http.MethodPost, "/product", strings.NewReader(`{
		"name": "Test product",
		"currency": "USD"
	}`))
	assert.Nil(t, err)

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getReq)
	assert.Equal(t, http.StatusOK, getRecorder.Code)

	var product models.Product
	err = json.Unmarshal(getRecorder.Body.Bytes(), &product)
	assert.Nil(t, err)

	currencyConverter := mockexchange.NewConverter()
	rate, err := currencyConverter.ConvertExchangeRate(fromCurrency, toCurrency)
	assert.Nil(t, err)

	assert.Equal(t, "Test product", product.Name)
	assert.Equal(t, initialPrice*rate, product.Price)
}

func TestGetProductWithUInvalidCurrency(t *testing.T) {
	router := run()

	req, err := http.NewRequest(http.MethodPost, "/products", strings.NewReader(`{
		"name": "Test product",
		"price": 100.0,
		"currency": "GBP"
	}`))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	getReq, err := http.NewRequest(http.MethodPost, "/product", strings.NewReader(`{
		"name": "Test product",
		"currency": "ABC"
	}`))
	assert.Nil(t, err)

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getReq)
	assert.Equal(t, http.StatusNotFound, getRecorder.Code)
}

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/keregocen/go-product-crud-api/pkg/mockexchange"
	"github.com/keregocen/go-product-crud-api/pkg/products"
	"github.com/keregocen/go-product-crud-api/pkg/storage"
)

func main() {
	router := run()

	log.Print("Starting the server on port 8080")
	if err := http.ListenAndServe(":8080", Middleware(router)); err != nil {
		log.Fatal(err)
	}
}

func run() *mux.Router {
	storageAPI := storage.NewStore()
	currencyConverterAPI := mockexchange.NewConverter()
	productHandler := products.Handler{StorageAPI: storageAPI, CurrencyConverterAPI: currencyConverterAPI}

	r := mux.NewRouter()
	r.HandleFunc("/products", productHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/products", productHandler.List).Methods(http.MethodGet)
	r.HandleFunc("/product", productHandler.Get).Methods(http.MethodPost)
	return r
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		h.ServeHTTP(w, r)
	})
}

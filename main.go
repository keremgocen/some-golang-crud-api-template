package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/keregocen/go-product-crud-api/pkg/mockexchange"
	"github.com/keregocen/go-product-crud-api/pkg/products"
	"github.com/keregocen/go-product-crud-api/pkg/storage"
)

func main() {
	ctxTimeoutSeconds := 5 * time.Second
	ctxReadHeaderTimeoutSeconds := 3 * time.Second

	storageAPI := storage.NewStore()
	currencyConverterAPI := mockexchange.NewConverter()
	service := products.NewService(storageAPI, currencyConverterAPI)

	r := mux.NewRouter()
	r.Handle("/products", middleware(http.HandlerFunc(service.Create))).Methods(http.MethodPost)
	r.Handle("/products", middleware(http.HandlerFunc(service.List))).Methods(http.MethodGet)
	r.Handle("/product", middleware(http.HandlerFunc(service.Get))).Methods(http.MethodPost)

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: ctxReadHeaderTimeoutSeconds,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error)
	go func() {
		log.Printf("Products service is starting.")
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-signalChan:
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeoutSeconds)
		defer cancel()
		log.Println("Server shutdown initiated.")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		h.ServeHTTP(w, r)
	})
}

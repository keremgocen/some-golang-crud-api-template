package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keregocen/go-product-crud-api/pkg/keyvalue"
	"github.com/keregocen/go-product-crud-api/pkg/storage"
)

const (
	ctxTimeout    = 5 * time.Second
	serverTimeout = 3 * time.Second
)

func main() {
	storageAPI := storage.NewStore()
	service := keyvalue.NewService(storageAPI)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.GET("/keyvalue/:name", service.GetKeyValue)
	router.POST("/keyvalue", service.PostKeyValue)

	httpServer := &http.Server{
		Addr:              ":5000",
		Handler:           router,
		ReadHeaderTimeout: serverTimeout,
		WriteTimeout:      serverTimeout,
		ReadTimeout:       serverTimeout,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error)
	go func() {
		log.Printf("keyvalue service is starting")
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-signalChan:
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		log.Println("server shutdown initiated")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}
}

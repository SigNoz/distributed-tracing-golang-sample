package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/NamanJain8/distributed-tracing-golang-sample/config"
	"github.com/NamanJain8/distributed-tracing-golang-sample/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

const serviceName = "payment-service"

var (
	srv        *http.Server
	paymentUrl string
	userUrl    string
)

func setupServer() {
	router := mux.NewRouter()
	router.HandleFunc("/payments/transfer/id/{userID}", transferAmount).Methods(http.MethodPut, http.MethodOptions)
	router.Use(utils.LoggingMW)
	router.Use(otelmux.Middleware(serviceName))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost},
	})

	srv = &http.Server{
		Addr:    paymentUrl,
		Handler: c.Handler(router),
	}

	log.Printf("Payment service running at: %s", paymentUrl)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to setup http server: %v", err)
	}
}

func main() {
	// read the config from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	paymentUrl = os.Getenv("PAYMENT_URL")
	userUrl = os.Getenv("USER_URL")

	// setup tracer
	tp := config.Init(serviceName)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go setupServer()

	<-sigint
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server shutdown failed")
	}
}

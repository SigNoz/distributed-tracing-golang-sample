package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/NamanJain8/distributed-tracing-golang-sample/config"
	"github.com/NamanJain8/distributed-tracing-golang-sample/datastore"
	"github.com/NamanJain8/distributed-tracing-golang-sample/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "order-service"

var (
	db       datastore.DB
	srv      *http.Server
	orderUrl string
	userUrl  string
	tracer   trace.Tracer
)

func setupServer() {
	router := mux.NewRouter()
	router.HandleFunc("/orders", createOrder).Methods(http.MethodPost, http.MethodOptions)
	router.Use(utils.LoggingMW)
	router.Use(otelmux.Middleware(serviceName))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost},
	})

	srv = &http.Server{
		Addr:    orderUrl,
		Handler: c.Handler(router),
	}

	log.Printf("Order service running at: %s", orderUrl)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to setup http server: %v", err)
	}
}

func initDB() {
	var err error
	if db, err = datastore.New(); err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
}

func main() {
	// read the config from .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}
	orderUrl = os.Getenv("ORDER_URL")
	userUrl = os.Getenv("USER_URL")

	// setup tracer
	tp := config.Init(serviceName)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	tracer = otel.Tracer(serviceName)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	initDB()
	go setupServer()

	<-sigint
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server shutdown failed")
	}
}

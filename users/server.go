package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/NamanJain8/distributed-tracing-golang-sample/datastore"
	"github.com/NamanJain8/distributed-tracing-golang-sample/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	db  datastore.DB
	srv *http.Server
)

func setupServer() {
	listenAddr := fmt.Sprintf("localhost:%d", httpPort())
	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/users/{userID}", getUser).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/users/{userID}", updateUser).Methods(http.MethodPut, http.MethodOptions)
	router.Use(utils.LoggingMW)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost},
	})

	srv = &http.Server{
		Addr:    listenAddr,
		Handler: c.Handler(router),
	}

	log.Printf("User service running on port: %d", httpPort())
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

func httpPort() int {
	p, ok := os.LookupEnv("USER_PORT")
	if !ok {
		return 8080
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatalf("incorrect port: %v", err)
	}

	return port
}

func main() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	initDB()
	go setupServer()

	<-sigint
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server shutdown failed")
	}
}

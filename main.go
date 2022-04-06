package main

import (
	"context"
	"log"
	"microservice/main.go/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)

	//create the handlers
	ph := handlers.NewProducts(l)

	//create a new server mux and register the handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	//create a new server
	s := &http.Server{
		Addr:         ":9090",           // configure the bind addres
		Handler:      sm,                // set the default handler
		IdleTimeout:  120 * time.Second, // max time to read request from the client
		ReadTimeout:  1 * time.Second,   // max time for connections using TCP keep-alive
		WriteTimeout: 1 * time.Second,   // max tie to write response to the client
	}

	//start the server
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recived terminate, graceful shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(ctx)
}

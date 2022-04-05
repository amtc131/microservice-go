package main

import (
	"context"
	"log"
	"microservice/main.go/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)

	//create the handlers
	//hh := handlers.NewHello(l)
	//	gh := handlers.NewGoodbye(l)
	ph := handlers.NewProducts(l)

	//create a new server mux and register the handlers
	sm := http.NewServeMux()
	//	sm.Handle("/hello", hh)
	//	sm.Handle("/goodbye", gh)
	sm.Handle("/", ph)

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

/*	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello world!!")
		d, _ := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Ooops", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(rw, "Hello  %s \n", d)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("GoodBye world!!")
	})
*/

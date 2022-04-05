package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//Hello is the simple handler
type Hello struct {
	l *log.Logger
}

//NewHello creates a new hello handler with the give logger
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

//ServeHTT implements the go http.Handler interface
//https://golang/pkg/net/http/·Handler
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello world!!")

	//read the body
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Ooops", http.StatusBadRequest)
		return
	}

	//write the response
	fmt.Fprintf(rw, "Hello  %s \n", d)

}

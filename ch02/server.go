package main

import (
	"fmt"
	"log"
	"net/http"
)

type myString string

func (s myString) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("myString=%s accessed", s)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "myString=%s", s)
	return
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/mystring/1", myString("1"))
	mux.Handle("/mystring/2", myString("2"))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

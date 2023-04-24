package main

import (
	"fmt"
	"log"
	"net/http"
)

var counter = 0

func homePage(w http.ResponseWriter, _ *http.Request) {
	counter += 1
	_, _ = fmt.Fprintln(w, "Counter: ", counter)
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Healthy!")
}

func ready(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Ready!")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/probe/liveness", healthcheck)
	http.HandleFunc("/probe/readiness", ready)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}

package main

import (
	"log"
	"net/http"
	"os"
)

var (
	port    string
	nodeURI string
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.Write([]byte("GOOD"))
	default:
		w.Write([]byte("Only POST method is supported."))
	}
}

func handleGetProof(w http.ResponseWriter, r *http.Request) {

}

func main() {
	var ok bool
	port, ok = os.LookupEnv("PORT")
	if !ok {
		port = "5000"
	}
	nodeURI, ok = os.LookupEnv("NODE_URI")
	if !ok {
		nodeURI = "tcp://localhost:26657"
	}
	http.HandleFunc("/request", handleRequest)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

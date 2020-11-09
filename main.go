package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	SetupRouter(r)
	port := os.Getenv("TSM_PORT")
	if port == "" {
		port = "8000"
	}
	log.Print("TSM | HTTP server listening on :", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

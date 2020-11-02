package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tsm/routers"
)

func main() {
	r := mux.NewRouter()
	routers.SetupRouter(r)
	log.Print("TSM | HTTP server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

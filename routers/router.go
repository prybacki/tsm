package routers

import (
	"github.com/gorilla/mux"
	"tsm/controllers"
)

func SetupRouter(r *mux.Router) {
	r.HandleFunc("/devices", controllers.HandlePost).Methods("POST")
}

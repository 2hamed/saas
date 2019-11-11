package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(d dispatcher) error {
	router := mux.NewRouter()

	router.HandleFunc("/new", NewJobHandler(d)).Methods("POST")
	router.HandleFunc("/result/{hash}", GetResultHandler(d)).Methods("GET")

	return http.ListenAndServe(":8080", router)
}

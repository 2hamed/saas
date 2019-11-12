package api

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func StartServer(d dispatcher) error {
	router := mux.NewRouter()

	router.HandleFunc("/api/new", NewJobHandler(d)).Methods("POST")
	router.HandleFunc("/api/result/{hash}", GetResultHandler(d)).Methods("GET")

	log.Infof("Started HTTP server on %s:%s", "", "8080")

	return http.ListenAndServe(":8080", router)
}

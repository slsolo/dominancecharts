package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{trait}", app.TraitGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{trait}/{name}", app.TraitGetValueHandler).Methods(http.MethodGet)

	return r
}

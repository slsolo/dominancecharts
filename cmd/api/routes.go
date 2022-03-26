package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/healthcheck", srv.healthcheckHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{trait}/names", srv.TraitNamesHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{trait}/", srv.TraitGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{trait}/{name}/", srv.TraitGetValueHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{trait}/compare/{first}/{second}", srv.TraitCompareHandler).Methods(http.MethodGet)

	return r
}

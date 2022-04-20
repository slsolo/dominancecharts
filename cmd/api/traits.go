package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"

	"github.com/slsolo/dominancecharts/internal/data"
)

func (application *application) TraitGetHandler(w http.ResponseWriter, r *http.Request) {
	var traits []data.Trait
	vars := mux.Vars(r)
	placed := application.models.Data
	mtraits := placed[strings.Title(vars["trait"])]
	for _, v := range mtraits {
		traits = append(traits, v)
	}

	if len(traits) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No Furs found.\n")
		return
	}
	sort.Slice(traits, func(i, j int) bool { return traits[i].Name < traits[j].Name })
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(traits)

}

func (app *application) TraitGetValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	placed := app.models.Data
	mtraits := placed[strings.Title(vars["trait"])]

	trait, ok := mtraits[vars["name"]]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s %s not found.\n", vars["trait"], vars["name"])
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trait)
}

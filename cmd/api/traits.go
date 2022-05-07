package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/slsolo/dominancecharts/internal/data"
)

func (app *application) TraitGetHandler(w http.ResponseWriter, r *http.Request) {
	var traits []data.Trait
	trait := chi.URLParam(r, "trait")
	placed := app.models.Data
	mtraits := placed[strings.Title(trait)]
	for _, v := range mtraits {
		traits = append(traits, v)
	}

	if len(traits) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No %ss found.\n", trait)
		return
	}
	sort.Slice(traits, func(i, j int) bool { return traits[i].Name < traits[j].Name })
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(traits)

}

func (app *application) TraitGetValueHandler(w http.ResponseWriter, r *http.Request) {
	traitType := chi.URLParam(r, "trait")
	traitName := chi.URLParam(r, "name")
	placed := app.models.Data
	mtraits := placed[strings.Title(traitType)]

	trait, ok := mtraits[traitName]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s %s not found.\n", traitType, traitName)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trait)
}

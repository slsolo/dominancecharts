package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (s *application) TraitNamesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	values, err := s.getData("names", strings.Title(vars["trait"]), "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	if len(values) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No %ss found.\n", vars["trait"])
		return
	}
	sort.Strings(values)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(values)
}

func (s *application) TraitGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	placed := make(map[string]int)
	values, err := s.getData("values", strings.Title(vars["trait"]), "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	if len(values) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No %ss found.\n", vars["trait"])
		return
	}

	respLen := len(values)
	for i := 0; i < respLen; i += 2 {
		placed[values[i]], _ = strconv.Atoi(values[i+1])
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(placed)
}

func (s *application) TraitGetValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	values, err := s.getData("names", strings.Title(vars["trait"]), vars["name"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	if len(values) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s %s not found.\n", vars["trait"], vars["name"])
		return
	}
	var i int
	i, err = strconv.Atoi(values[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "value (%s) is not an integer\n", values[0])
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}

func (s *application) TraitCompareHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trait := strings.Title(vars["trait"])
	first := vars["first"]
	second := vars["second"]
	var fi, si int
	values, err := s.getData("single", trait, first)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	if len(values) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s %s not found.\n", trait, first)
		return
	}

	fi, err = strconv.Atoi(values[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "value (%s) is not an integer\n", values[0])
		return
	}

	values, err = s.getData("single", trait, second)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	if len(values) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s %s not found.\n", trait, second)
		return
	}

	si, err = strconv.Atoi(values[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "value (%s) is not an integer\n", values[0])
		return
	}

	if first == second {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Both %ss are the same. Please change your selection and try again.", trait)
	} else {
		w.WriteHeader(http.StatusOK)
		if fi < si {
			fmt.Fprintf(w, "%s is dominant to %s", first, second)
		} else {
			fmt.Fprintf(w, "%s is recessive to %s", first, second)
		}

	}
}

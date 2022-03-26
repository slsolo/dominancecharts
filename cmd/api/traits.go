package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

func (s *application) TraitNamesHandler(w http.ResponseWriter, r *http.Request) {
	values, err := s.getData("names", strings.Title(c.Param("trait")), "")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(values) == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No %ss found.", c.Param("trait")))
		return
	}
	sort.Strings(values)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(values).Encode(w)
}

func (s *application) TraitGetHandler(w http.ResponseWriter, r *http.Request) {
	placed := make(map[string]int)
	values, err := s.getData("values", strings.Title(c.Param("trait")), "")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(values) == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("no %ss found", c.Param("trait")))
		return
	}

	respLen := len(values)
	for i := 0; i < respLen; i += 2 {
		placed[values[i]], _ = strconv.Atoi(values[i+1])
	}
	c.JSON(http.StatusOK, placed)
}

func (s *application) TraitGetValueHandler(w http.ResponseWriter, r *http.Request) {
	values, err := s.getData("names", strings.Title(c.Param("trait")), c.Param("name"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(values) == 0 {
		c.String(http.StatusNotFound, "%s %s not found\n", c.Param("Trait"), c.Param("name"))
		return
	}
	var i int
	i, err = strconv.Atoi(values[0])
	if err != nil {
		c.String(http.StatusInternalServerError, "Value is not an integer")
		return
	}
	c.JSON(200, i)
}

func (s *application) TraitCompareHandler(w http.ResponseWriter, r *http.Request) {
	trait := strings.Title(c.Param("trait"))
	first := c.Param("first")
	second := c.Param("second")
	var fi, si int
	values, err := s.getData("single", trait, first)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(values) == 0 {
		c.String(http.StatusNotFound, "%s %s not found\n", trait, first)
		return
	}

	fi, err = strconv.Atoi(values[0])
	if err != nil {
		log.Printf("trying to parse %s\n", values[0])
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	values, err = s.getData("single", trait, second)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(values) == 0 {
		c.String(http.StatusNotFound, "%s %s not found\n", trait, second)
		return
	}

	si, err = strconv.Atoi(values[0])
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if first == second {
		c.String(400, "Both %ss are the same. Please change your selection and try again.", trait)
	} else if fi < si {
		c.String(200, fmt.Sprintf("%s is dominant to %s", first, second))
	} else {
		c.String(200, fmt.Sprintf("%s is recessive to %s", first, second))
	}

}

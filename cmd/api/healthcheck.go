package main

import (
	"fmt"
	"net/http"
)

func (srv *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
}

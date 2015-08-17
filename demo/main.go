// chris 081715

//go:generate govalid valid.v

// Demo server for validation suite.
//
// Listens on /user/create and /group/create, and validates separate
// input types for each defined in valid.v.
package main

import (
	"flag"
	"io"
	"log"

	"net/http"
)

func handleError(w http.ResponseWriter, status int, err error) {
	w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(status)
	io.WriteString(w, err.Error())
}

func handleUser(w http.ResponseWriter, req *http.Request) {
	ud, err := getReqInput(req)
	if err != nil {
		handleError(w, 500, err)
		return
	}
	ui, err2 := validateUserInput(ud)
	if err2 != nil {
		handleError(w, 500, err2)
		return
	}
	password, err3 := hashPass(ui.password)
	if err3 != nil {
		handleError(w, 500, err3)
		return
	}
	io.WriteString(w, "ok")
	ui.password = ""
	log.Println("user input", ui)
	log.Println("bcrypt'd user password", string(password))
}

func handleGroup(w http.ResponseWriter, req *http.Request) {
	gd, err := getReqInput(req)
	if err != nil {
		handleError(w, 500, err)
		return
	}
	gi, err2 := validateGroupInput(gd)
	if err2 != nil {
		handleError(w, 500, err2)
		return
	}
	io.WriteString(w, "ok")
	log.Println("group input", gi)
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	http.HandleFunc("/user/create", handleUser)
	http.HandleFunc("/group/create", handleGroup)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

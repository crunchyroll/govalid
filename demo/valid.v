// chris 081715

package main

import (
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func getReqInput(req *http.Request) (map[string]string, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}
	v := req.Form
	r := make(map[string]string, len(v))
	for k, vs := range v {
		r[k] = vs[0]
	}
	return r, nil
}

func hashPass(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

type userInput struct {
	username string `valid:"max:16"`
	password []byte `valid:"max:128"`

	fname string `valid:"max:128,def"`
	lname string `valid:"max:128,def"`

	email *mail.Address `valid:"max:1024"`

	age uint `valid:"min:13,max:200,def"`
}

type groupInput struct {
	name string `valid:"max:16"`
	descr string `valid:"max:1024,def"`
}

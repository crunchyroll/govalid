// Test .v for adding bounded string and number suport.

package fgsfds

import (
	"net/mail"
	"net/url"
)

type strings struct {
	short string `valid:"max:3"`
	min   string `valid:"min:3"`
	band  string `valid:"min:4,max:10"`

	optmax string `valid:"max:4,def"`

	def string `valid:"def"`
	def2 string `valid:"def:\"default\""`
}

type numbers struct {
	small int `valid:"max:100"`
	min   int `valid:"min:20"`
	band  int `valid:"max:256,min:20"`

	neg int `valid:"min:-10,max:-5"`

	def int32 `valid:"def"`
	def2 int32 `valid:"def:-3"`

	f float32 `valid:"def:3.14,min:-100.123"`
}

type errs struct {
	// These would cause an error if the field name temporary
	// variables were not namespaced.
	ok  int
	err float32
}

type net struct {
	email *mail.Address `valid:"def"`

	homepage *url.URL `valid:"max:2048"`
}

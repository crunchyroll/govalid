// Test .v for adding *url.URL and *mail.Address support.

package whatever

import (
	"net/mail"
	"net/url"
)

type userInput struct {
	username string
	password string
	age      uint8
	email    *mail.Address
	homepage *url.URL
}

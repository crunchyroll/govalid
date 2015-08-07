// Test .v for adding *url.URL and *mail.Address support.

// TODO Come up with reasonable functionality (or document) when you are
// missing the url or mail imports.

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

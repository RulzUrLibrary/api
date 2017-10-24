package utils

import (
	"errors"
)

var (
	ErrInvalidIsbn     = errors.New("isbn is invalid")
	ErrNotFound        = errors.New("entity not found")
	ErrExists          = errors.New("entity already exists")
	ErrPageNotFound    = errors.New("page not found")
	ErrHTMLHandler     = errors.New("endpoint only available for text/html requests")
	ErrNoProduct       = errors.New("no product corresponding to isbn")
	ErrParsingProduct  = errors.New("error when parsing product page")
	ErrUserAuth        = errors.New("username or password incorrect")
	ErrUserExists      = errors.New("user already registered")
	ErrNotUser         = errors.New("you need to provide auth")
	ErrCaptcha         = errors.New("captcha hitted")
	ErrGoogleAuth      = errors.New("google auth failed")
	ErrOffset          = errors.New("offset must be greater than 0")
	ErrLimit           = errors.New("limit must be between 0 and 50")
	ErrAlreadyActivate = errors.New("account already activated")
)

package auth

import (
	"github.com/labstack/echo"
	"github.com/paul-bismuth/library/ext/db"
	"github.com/paul-bismuth/library/ext/google"
	"github.com/paul-bismuth/library/utils"
	"strings"
)

type Auth struct {
	cache  map[string]utils.User
	DB     *db.DB
	Logger echo.Logger
}

func New(l echo.Logger, d *db.DB) *Auth {
	return &Auth{make(map[string]utils.User), d, l}
}

func (auth *Auth) Login(username, password string) (user utils.User, err error) {
	auth.Logger.Infof("user %s is attempting to log in", username)

	if user, ok := auth.cache[username+password]; ok {
		auth.Logger.Infof("using cache on user: %+v", user)
		return user, nil
	}
	suffix := strings.TrimLeftFunc(username, func(r rune) bool { return r != '@' })

	switch suffix {
	case "@gmail.com":
		user, err = google.Auth(auth.DB, username, password)
	default:
		user, err = auth.DB.Auth(username, password)
	}
	if err != nil {
		return
	}
	auth.Logger.Infof("caching user: %+v", user)
	auth.cache[username+password] = user
	return
}

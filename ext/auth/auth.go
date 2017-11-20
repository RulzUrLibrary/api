package auth

import (
	"github.com/labstack/echo"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/ext/google"
	"github.com/rulzurlibrary/api/utils"
)

type Cache interface {
	Get(string) (*utils.User, bool)
	Set(string, *utils.User)
}

type DefaultCache struct {
	cache  map[string]*utils.User
	logger echo.Logger
}

func NewDefaultCache(logger echo.Logger) DefaultCache {
	return DefaultCache{map[string]*utils.User{}, logger}
}

func (df DefaultCache) Get(key string) (user *utils.User, ok bool) {
	user, ok = df.cache[key]
	if ok {
		df.logger.Infof("using cache on user: %+v", user)
	}
	return
}

func (df DefaultCache) Set(key string, value *utils.User) {
	df.logger.Infof("caching user: %+v", value)
	df.cache[key] = value
}

type Auth struct {
	cache  Cache
	DB     *db.DB
	Logger echo.Logger
}

func New(l echo.Logger, d *db.DB, c Cache) *Auth {
	return &Auth{c, d, l}
}

func (auth *Auth) Login(username, password string) (user *utils.User, err error) {
	auth.Logger.Infof("user %s is attempting to log in", username)

	if user, ok := auth.cache.Get(username + password); ok {
		return user, nil
	}

	switch suffix := utils.MailAddress(username); suffix {
	case "@gmail.com":
		user, err = google.Auth(auth.DB, username, password)
	default:
		user, err = auth.DB.Auth(username, password)
	}
	if err != nil {
		return
	}
	auth.cache.Set(username+password, user)
	return
}

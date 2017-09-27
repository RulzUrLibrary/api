package app

import (
	"github.com/ixday/echo-hello/ext/google"
	"github.com/ixday/echo-hello/utils"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

var cache = make(map[string]utils.User) // some in memory cache

func (app *Application) auth(username, password string) (utils.User, error) {
	suffix := strings.TrimLeftFunc(username, func(r rune) bool { return r != '@' })
	switch suffix {
	case "@gmail.com":
		return google.Auth(app.Database, username, password)
	}
	return app.Database.Auth(username, password)
}

func unauthorized(c echo.Context, err error) error {
	c.Response().Header().Set(echo.HeaderWWWAuthenticate, "Basic realm=Restricted")
	return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
}

func save(c echo.Context, next echo.HandlerFunc, user utils.User) error {
	c.Set("user", user)
	return next(c)
}

func (app *Application) BasicAuth(strict bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username, password, ok := c.Request().BasicAuth()
			if ok {
			} else if strict {
				return unauthorized(c, utils.ErrNotUser)
			} else {
				return next(c)
			}
			app.Logger.Infof("user %s is attempting to log in", username)

			if user, ok := cache[username+password]; ok {
				app.Logger.Infof("using cache on user: %+v", user)
				return save(c, next, user)
			}
			user, err := app.auth(username, password)
			if err == nil {
			} else if err == utils.ErrUserAuth {
				return unauthorized(c, err)
			} else {
				return err
			}
			app.Logger.Infof("caching user: %+v", user)
			cache[username+password] = user
			return save(c, next, user)
		}
	}
}

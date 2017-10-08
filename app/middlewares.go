package app

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/paul-bismuth/library/utils"
	"net/http"
	"net/url"
)

const (
	KEY_SIZE     = 16
	SESSION_KEY  = "session"
	SESSION_NAME = "rulz"
)

func ContentType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		if req.Method == echo.GET {
			return next(c)
		}
		ct := req.Header.Get(echo.HeaderContentType)
		if ct != echo.MIMEApplicationJSON {
			return echo.NewHTTPError(
				http.StatusBadRequest, "API only support application/json content type",
			)
		}
		return next(c)
	}
}

func unauthorized(c echo.Context, err error) error {
	c.Response().Header().Set(echo.HeaderWWWAuthenticate, "Basic realm=Restricted")
	return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
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

			switch user, err := app.Auth.Login(username, password); err {
			case nil:
				c.Set("user", user)
			case utils.ErrUserAuth:
				return unauthorized(c, err)
			default:
				return err
			}
			return next(c)
		}
	}
}

func Protected(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if user := c.Get("user"); user != nil {
			return next(c)
		}

		v := url.Values{}
		v.Set("next", c.Request().URL.String())

		u, err := url.Parse(c.Echo().Reverse("auth"))
		if err != nil {
			return err
		}
		u.RawQuery = v.Encode()
		return c.Redirect(http.StatusSeeOther, u.String())
	}
}

func CookieAuth(dev bool) echo.MiddlewareFunc {
	store := sessions.NewCookieStore([]byte{
		88, 111, 50, 219, 142, 197, 174, 166, 14, 229, 175, 140, 165, 97, 112, 62,
		100, 74, 227, 150, 198, 247, 19, 76, 90, 160, 247, 44, 100, 200, 25, 163,
	})
	if !dev {
		store = sessions.NewCookieStore(securecookie.GenerateRandomKey(KEY_SIZE))
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			session, err := store.Get(request, SESSION_NAME)
			if err != nil {
				// reinit cookies
				session, _ = store.New(request, SESSION_NAME)
			} else {
				if user, ok := session.Values["user"]; ok {
					c.Set("user", user)
				}
			}
			c.Set(SESSION_KEY, session)
			return next(c)
		}
	}
}

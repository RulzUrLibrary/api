package app

import (
	"github.com/gorilla/sessions"
	"github.com/paul-bismuth/library/utils"
	"net/http"
)

func WEBUserGet(c *Context) error {
	return c.Render(http.StatusOK, "user.html", map[string]interface{}{
		"user": c.Get("user"),
	})
}

func WEBUserLogout(c *Context) error {
	session, _ := c.Get("session").(*sessions.Session)
	session.Values["user"] = nil
	session.AddFlash("you have been successfully logged out!")
	c.Set("user", nil)
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("index"))
}

func WEBAuthGet(c *Context) error {
	return c.Render(http.StatusOK, "auth.html", map[string]interface{}{})
}

func WEBAuthPost(c *Context) error {
	creds := struct {
		User     string `form:"user"`
		Password string `form:"password"`
		Token    string `form:"token"`
		Next     string `query:"next"`
	}{Next: c.Echo().Reverse("books")}

	if err := c.Bind(&creds); err != nil {
		return err
	}
	user, err := c.Auth.Login(creds.User, utils.DefaultS(creds.Token, creds.Password))
	if err != nil {
		return c.Render(http.StatusUnauthorized, "auth.html", map[string]interface{}{
			"error": err,
		})
	}
	session, _ := c.Get("session").(*sessions.Session)
	session.Values["user"] = &user
	c.Set("user", &user)
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, creds.Next)
}

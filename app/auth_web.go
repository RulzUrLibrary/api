package app

import (
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBUserGet(c *Context) error {
	return c.Render(http.StatusOK, "user.html", map[string]interface{}{
		"user": c.Get("user"),
	})
}

func WEBUserLogout(c *Context) error {
	flash := utils.Flash{utils.FlashSuccess, "you have been successfully logged out!"}
	if err := c.SaveUser(nil, flash); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("index"))
}

func WEBUserNewGet(c *Context) error {
	return c.Render(http.StatusOK, "new.html", map[string]interface{}{})
}

func WEBUserNewPost(c *Context) error {
	badRequest := func(msg string) error {
		return c.Render(http.StatusBadRequest, "new.html", map[string]interface{}{
			"error": msg,
		})
	}
	creds := struct {
		User         string `form:"user"`
		Password     string `form:"password"`
		Confirmation string `form:"confirmation"`
	}{}

	if err := c.Bind(&creds); err != nil {
		return badRequest(err.Error())
	}
	if creds.Confirmation != creds.Password {
		return badRequest("passwords don't match")
	}
	if len(creds.Password) < 2 {
		return badRequest("user password must be at least 8 characters long")
	}
	switch suffix := utils.MailAddress(creds.User); suffix {
	case "@gmail.com":
		return badRequest("username should not be from gmail, signin directly if you want to use it")
	}
	user, err := c.DB.NewUser(creds.User, creds.Password)
	switch err {
	case nil:
		flash := utils.Flash{utils.FlashSuccess, "welcome to rulz!"}
		if err := c.SaveUser(user, flash); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("books"))
	case utils.ErrUserExists:
		return badRequest(err.Error())
	}
	return err
}

func WEBAuthGet(c *Context) error {
	return c.Render(http.StatusOK, "auth.html", map[string]interface{}{})
}

func WEBAuthPost(c *Context) error {
	creds := struct {
		User     string `form:"user"`
		Password string `form:"password"`
		Token    string `form:"token"`
		Next     string `form:"next"`
	}{}

	if err := c.Bind(&creds); err != nil {
		return err
	}
	user, err := c.Auth.Login(creds.User, utils.DefaultS(creds.Token, creds.Password))
	if err != nil {
		return c.Render(http.StatusUnauthorized, "auth.html", map[string]interface{}{
			"error": err,
		})
	}
	if err := c.SaveUser(user); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, utils.DefaultS(creds.Next, c.Echo().Reverse("books")))
}

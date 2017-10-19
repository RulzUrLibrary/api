package app

import (
	"github.com/rulzurlibrary/api/ext/validator"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

type dict = utils.Dict

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

func WEBUserPassword(c *Context) error {
	return nil
}

func WEBUserReset(c *Context) error {
	return nil
}

func WEBUserNewGet(c *Context) error {
	query := struct {
		Email        string `query:"email" validate:"email,gmail"`
		Password     string
		Confirmation string
	}{}
	errs := dict{}

	if err := c.Bind(&query); err != nil {
		return err
	}
	if err := c.Validate(&query); err != nil {
		errs = validator.Dump(err, map[string]map[string]string{
			"email": map[string]string{
				"email": "invalid email address",
				"gmail": "email should not be from gmail, signin directly if you want to use it",
			},
		})
	}
	return c.Render(http.StatusOK, "new.html", dict{"error": errs, "form": query})
}

func WEBUserNewPost(c *Context) error {
	creds := struct {
		Email        string `form:"email" validate:"required,email,gmail"`
		Password     string `form:"password" validate:"required,gt=8,eqfield=Confirmation"`
		Confirmation string `form:"confirmation"`
	}{}
	badRequest := func(err interface{}) error {
		return c.Render(http.StatusBadRequest, "new.html", dict{"error": err, "form": creds})
	}

	if err := c.Bind(&creds); err != nil {
		return err
	}
	if err := c.Validate(&creds); err != nil {
		return badRequest(validator.Dump(err, map[string]map[string]string{
			"email": map[string]string{
				"required": "email is required",
				"email":    "invalid email address",
				"gmail":    "email should not be from gmail, signin directly if you want to use it",
			},
			"password": map[string]string{
				"required": "password is required",
				"gt":       "password must be at least 8 characters long",
				"eqfield":  "passwords must match",
			},
		}))
	}

	user, err := c.DB.NewUser(creds.Email, creds.Password)
	switch err {
	case nil:
		flash := utils.Flash{utils.FlashSuccess, "welcome to rulz!"}
		if err := c.SaveUser(user, flash); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("books"))
	case utils.ErrUserExists:
		return badRequest(dict{"email": err.Error()})
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

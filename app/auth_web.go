package app

import (
	"github.com/rulzurlibrary/api/ext/validator"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBUserGet(c *Context) error {
	user := c.Get("user").(*utils.User)
	return c.Render(http.StatusOK, "user.html",
		dict{"error": dict{}, "user": user, "form": struct {
			Old  string
			New  string
			Conf string
		}{}, "misc": struct{ Valid bool }{utils.ValidMailProvider(user.Email)}},
	)
}

func WEBUserLang(c *Context) error {
	c.SetCookie(&http.Cookie{
		Name: "lang", Value: c.FormValue("lang"), HttpOnly: false, Path: "/",
	})
	return c.Redirect(http.StatusSeeOther, c.FormValue("next"))
}

func WEBUserLogout(c *Context) error {
	flash := utils.Flash{utils.FlashSuccess, utils.FLASH_LOGOUT}
	if err := c.SaveUser(nil, flash); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("index"))
}

func WEBUserResetPost(c *Context) error {
	creds := struct {
		Old  string `form:"old" validate:"required"`
		New  string `form:"new" validate:"required,gt=8,eqfield=Conf,nefield=Old"`
		Conf string `form:"confirmation"`
	}{}
	user := c.Get("user").(*utils.User)
	badRequest := func(err interface{}) error {
		return c.Render(http.StatusBadRequest, "user.html", dict{
			"error": err, "user": user, "form": creds,
			"misc": struct{ Valid bool }{utils.ValidMailProvider(user.Email)},
		})
	}
	if err := c.Bind(&creds); err != nil {
		return err
	}
	if err := c.Validate(&creds); err != nil {
		return badRequest(validator.Dump(err, map[string]dict{
			"old": dict{"required": utils.OLD_PASSWORD_REQUIRED},
			"new": dict{
				"required": utils.PASSWORD_REQUIRED, "gt": utils.PASSWORD_LEN,
				"eqfield": utils.PASSWORD_EQFIELD, "nefield": utils.PASSWORD_NEQFIELD,
			},
		}))
	}
	if count, err := c.DB.ChangePassword(
		creds.New, creds.Old, user.Id,
	); err != nil {
		return err
	} else if count == 0 {
		return badRequest(dict{"old": utils.PASSWORD_INVALID})
	}
	if err := c.Flashes(utils.Flash{utils.FlashSuccess, utils.FLASH_PASSWORD}); err != nil {
		return err
	}
	return WEBUserGet(c)
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
		errs = validator.Dump(err, map[string]dict{
			"email": dict{"email": utils.EMAIL_INVALID, "gmail": utils.EMAIL_GMAIL},
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
		return badRequest(validator.Dump(err, map[string]dict{
			"email": dict{
				"required": utils.EMAIL_REQUIRED, "email": utils.EMAIL_INVALID,
				"gmail": utils.EMAIL_GMAIL,
			},
			"password": dict{
				"required": utils.PASSWORD_REQUIRED, "gt": utils.PASSWORD_LEN,
				"eqfield": utils.PASSWORD_EQFIELD,
			},
		}))
	}

	user, _, err := c.DB.NewUser(creds.Email, creds.Password)
	switch err {
	case nil:
		flash := utils.Flash{utils.FlashSuccess, utils.FLASH_WELCOME}
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
	return c.Render(http.StatusOK, "auth.html", dict{
		"error": dict{}, "form": struct {
			Email    string
			Password string
		}{},
	})
}

func WEBAuthPost(c *Context) error {
	creds := struct {
		Email    string `form:"email" validate:"required"`
		Password string `form:"password" validate:"required"`
		Token    string `form:"token"`
		Next     string `form:"next"`
	}{}
	render := func(code int, errs dict) error {
		return c.Render(code, "auth.html", dict{"error": errs, "form": creds})
	}

	if err := c.Bind(&creds); err != nil {
		return err
	}
	creds.Password = utils.DefaultS(creds.Token, creds.Password)
	if err := c.Validate(&creds); err != nil {
		return render(http.StatusBadRequest, validator.Dump(err, map[string]dict{
			"email":    dict{"required": utils.EMAIL_REQUIRED},
			"password": dict{"required": utils.PASSWORD_REQUIRED},
		}))
	}
	user, err := c.Auth.Login(creds.Email, creds.Password)
	if err != nil {
		return render(http.StatusUnauthorized, dict{"auth": utils.AUTH_ERR})
	}
	if err := c.SaveUser(user); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, utils.DefaultS(creds.Next, c.Echo().Reverse("books")))
}

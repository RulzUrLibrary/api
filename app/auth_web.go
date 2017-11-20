package app

import (
	"github.com/RulzUrLibrary/api/ext/validator"
	"github.com/RulzUrLibrary/api/utils"
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

func WEBUserChange(c *Context) error {
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
	if count, err := c.App.Database.PasswordChange(
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

func WEBUserResetGet(c *Context) error {
	email := ""
	if user, ok := c.Get("user").(*utils.User); ok {
		email = user.Email
	}
	return c.Render(http.StatusOK, "reset.html", dict{"email": email, "error": dict{}})
}

func WEBUserResetPost(c *Context) error {
	query := struct {
		Email string `form:"email" validate:"required,email,gmail"`
	}{}

	if err := c.Bind(&query); err != nil {
		return err
	}
	if err := c.Validate(&query); err != nil {
		errs := validator.Dump(err, map[string]dict{
			"email": dict{
				"email": utils.EMAIL_INVALID, "gmail": utils.EMAIL_GMAIL_RESET,
				"required": utils.EMAIL_REQUIRED,
			},
		})
		return c.Render(http.StatusBadRequest, "reset.html", dict{"email": query.Email, "error": errs})
	}

	if link, err := c.App.Database.CreateReset(query.Email); err != nil {
		return err
	} else if link != "" {
		link = c.ReverseAbs("reinit", link)
		if err := c.App.Smtp.ResetMail(c, query.Email, link); err != nil {
			return err
		}
	}

	return c.Render(http.StatusOK, "sent.html", dict{"email": query.Email})
}

func WEBUserReinit(c *Context) error {
	reset := c.Param("id")
	form := struct {
		Password     string `form:"password" validate:"required,gt=8,eqfield=Confirmation"`
		Confirmation string `form:"confirmation"`
	}{}
	errors := dict{}
	render := func(code int) error {
		return c.Render(code, "reinit.html", dict{"form": form, "errors": errors})
	}

	if err := dynamic(c, "users", "reset", reset); err != nil {
		return err
	}

	if c.Request().Method == http.MethodPost {
		if err := c.Bind(&form); err != nil {
			return err
		}
		if err := c.Validate(&form); err != nil {
			errors = validator.Dump(err, map[string]dict{
				"password": dict{
					"required": utils.PASSWORD_REQUIRED, "gt": utils.PASSWORD_LEN,
					"eqfield": utils.PASSWORD_EQFIELD,
				},
			})
			return render(http.StatusBadRequest)
		}
		if err := c.App.Database.PasswordReset(form.Password, reset); err != nil {
			return err
		}
		return c.RedirectWithFlash(utils.FLASH_PASSWORD)
	}

	return render(http.StatusOK)
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

	user, activate, err := c.App.Database.NewUser(creds.Email, creds.Password)
	switch err {
	case nil:
	case utils.ErrUserExists:
		return badRequest(dict{"email": err.Error()})
	default:
		return err
	}
	flash := utils.Flash{utils.FlashSuccess, utils.FLASH_WELCOME}
	if err := c.SaveUser(user, flash); err != nil {
		c.App.Database.MustDeleteUser(user)
		return err
	}
	activate = c.ReverseAbs("activate", activate)
	if err := c.App.Smtp.ActivationMail(c.Context, creds.Email, activate); err != nil {
		c.App.Database.MustDeleteUser(user)
		return err
	}

	return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("books"))
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
	user, err := c.App.Auth.Login(creds.Email, creds.Password)
	if err != nil {
		return render(http.StatusUnauthorized, dict{"auth": utils.ERR_AUTH})
	}
	if err := c.SaveUser(user); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, utils.DefaultS(creds.Next, c.Echo().Reverse("books")))
}

func WEBUserActivate(c *Context) error {
	activate := c.Param("id")
	if err := dynamic(c, "users", "activate", activate); err != nil {
		return err
	}
	switch err := c.App.Database.DeleteActivate(activate); err {
	case nil:
	case utils.ErrAlreadyActivate:
		return c.Render(http.StatusBadRequest, "error.html",
			dict{"code": http.StatusBadRequest, "msg": utils.ERR_ALREADY_ACTIVATED})
	default:
		return err
	}

	return c.RedirectWithFlash(utils.FLASH_ACTIVATED)
}

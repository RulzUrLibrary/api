package app

import (
	"github.com/rulzurlibrary/api/ext/validator"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

type dict = utils.Dict
type dictS = map[string]string

func WEBUserGet(c *Context) error {
	return c.Render(http.StatusOK, "user.html", map[string]interface{}{
		"user": c.Get("user"),
	})
}

func WEBUserLogout(c *Context) error {
	flash := utils.Flash{utils.FlashSuccess, utils.FLASH_LOGOUT}
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
		errs = validator.Dump(err, map[string]dictS{
			"email": dictS{"email": utils.EMAIL_INVALID, "gmail": utils.EMAIL_GMAIL},
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
		return badRequest(validator.Dump(err, map[string]dictS{
			"email": dictS{
				"required": utils.EMAIL_REQUIRED, "email": utils.EMAIL_INVALID,
				"gmail": utils.EMAIL_GMAIL,
			},
			"password": dictS{
				"required": utils.PASSWORD_REQUIRED, "gt": utils.PASSWORD_LEN,
				"eqfield": utils.PASSWORD_EQFIELD,
			},
		}))
	}

	user, err := c.DB.NewUser(creds.Email, creds.Password)
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
		return render(http.StatusBadRequest, validator.Dump(err, map[string]dictS{
			"email":    dictS{"required": utils.EMAIL_REQUIRED},
			"password": dictS{"required": utils.PASSWORD_REQUIRED},
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

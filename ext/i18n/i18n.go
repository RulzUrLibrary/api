package i18n

import (
	"github.com/labstack/echo"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/rulzurlibrary/api/utils"
	"path/filepath"
)

type Configuration struct {
	Path    string
	Default string
}

type I18n struct {
	Configuration
}

func GetI18n(c echo.Context) i18n.TranslateFunc {
	return c.Get("t").(i18n.TranslateFunc)
}

func (i *I18n) SetI18n(c echo.Context) {
	cookieLang := ""
	if cookie, err := c.Cookie("lang"); err == nil {
		cookieLang = cookie.Value
	}
	acceptLang := c.Request().Header.Get("Accept-Language")
	defaultLang := i.Default // known valid language

	c.Set("t", i18n.MustTfunc(cookieLang, acceptLang, defaultLang))
	c.Set("lang", utils.DefaultS(cookieLang, acceptLang, defaultLang)[0:2])
}

func New(config Configuration) *I18n {
	matches, _ := filepath.Glob(filepath.Join(config.Path, "*.all.json"))
	for _, match := range matches {
		i18n.MustLoadTranslationFile(match)
	}
	return &I18n{config}
}

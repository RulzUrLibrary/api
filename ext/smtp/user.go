package smtp

import (
	"bytes"
	"github.com/labstack/echo"
	"github.com/RulzUrLibrary/api/ext/i18n"
	"github.com/RulzUrLibrary/api/utils"
)

func (s *Smtp) mail(c echo.Context, subject, dest, tplt string, data utils.Dict) error {
	subject = i18n.GetI18n(c)(subject)
	renderer := c.Echo().Renderer
	buffer := bytes.Buffer{}
	err := renderer.Render(&buffer, tplt, data, c)
	if err != nil {
		return err
	}
	return s.NewMail("noreply").
		To("", dest).
		Subject(subject).
		Body(buffer.Bytes()).
		Send()
}

func (s *Smtp) ActivationMail(c echo.Context, dest, link string) error {
	return s.mail(c, utils.SUBJECT_ACTIVATION, dest, "mail/activation.html",
		utils.Dict{"activate": link})
}

func (s *Smtp) ResetMail(c echo.Context, dest, link string) error {
	return s.mail(c, utils.SUBJECT_RESET, dest, "mail/reset.html",
		utils.Dict{"reset": link})
}

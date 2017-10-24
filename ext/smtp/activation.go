package smtp

import (
	"bytes"
	"github.com/labstack/echo"
	"github.com/rulzurlibrary/api/ext/i18n"
	"github.com/rulzurlibrary/api/utils"
)

func (s *Smtp) ActivationMail(c echo.Context, dest, link string) error {
	renderer := c.Echo().Renderer
	buffer := bytes.Buffer{}
	subject := i18n.GetI18n(c)(utils.ACTIVATION_SUBJECT)
	err := renderer.Render(&buffer, "activation.html", utils.Dict{"activate": link}, c)
	if err != nil {
		return err
	}
	return s.NewMail("noreply").
		To("", dest).
		Subject(subject).
		Body(buffer.Bytes()).
		Send()
}

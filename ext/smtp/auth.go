package smtp

import (
	"bytes"
	"github.com/labstack/echo"
	"github.com/rulzurlibrary/api/utils"
)

func (s *Smtp) ActivationMail(c echo.Context, dest, link string) error {
	renderer := c.Echo().Renderer
	buffer := bytes.Buffer{}
	err := renderer.Render(buffer, "activation.html", utils.Dict{"activate": activate}, c)

	if err != nil {
		return err
	}
	return s.NewMail("noreply").
		To("", dest).
		Subject().
		Body(buffer.String()).
		Send()
}

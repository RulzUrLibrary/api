package view

import (
	"github.com/CloudyKit/jet"
	"github.com/labstack/echo"
	"io"
)

type Configuration struct {
	Path        string
	Development bool
}

type View struct {
	*jet.Set
}

func (v *View) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tplt, err := v.GetTemplate(name)
	if err != nil {
		// template could not be loaded
		return err
	}
	return tplt.Execute(w, nil, data)
}

func New(config Configuration) *View {
	view := &View{jet.NewHTMLSet(config.Path)}
	view.SetDevelopmentMode(config.Development)
	return view
}

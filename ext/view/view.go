package view

import (
	"github.com/CloudyKit/jet"
	"github.com/labstack/echo"
	"io"
	"reflect"
)

type Configuration struct {
	Path        string
	Development bool
	App         *echo.Echo
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

	// add a reverse url helper
	view.AddGlobalFunc("url", func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("url", 1, -1)
		name := a.Get(0).String()
		args := []interface{}{}
		for i := 1; i < a.NumOfArguments(); i++ {
			args = append(args, a.Get(i))
		}
		url := config.App.Reverse(name, args...)
		if url == "" {
			a.Panicf("route not found")
		}
		return reflect.ValueOf(url)
	})
	return view
}

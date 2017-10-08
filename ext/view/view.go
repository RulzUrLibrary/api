package view

import (
	"github.com/CloudyKit/jet"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"io"
	"reflect"
	"strings"
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
	vars := jet.VarMap{}
	tplt, err := v.GetTemplate(name)
	if err != nil {
		// template could not be loaded
		return err
	}
	session, _ := c.Get("session").(*sessions.Session)
	flashes := session.Flashes()
	if len(flashes) > 0 {
		session.Save(c.Request(), c.Response())
	}
	vars.Set("context", c)
	vars.Set("flashes", flashes)
	return tplt.Execute(w, vars, data)
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
	view.AddGlobalFunc("title", func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("title", 1, 1)
		s := a.Get(0).String()
		return reflect.ValueOf(strings.Title(s))
	})
	view.AddGlobalFunc("capitalize", func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("capitalize", 1, 1)
		s := a.Get(0).String()
		if len(s) > 0 {
			s = strings.ToUpper(s[:1]) + s[1:]
		}
		return reflect.ValueOf(s)
	})
	view.AddGlobalFunc("valid", func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("valid", 1, 1)
		if a.Get(0).IsNil() {
			return reflect.ValueOf(false)
		} else {
			return reflect.ValueOf(true)
		}
	})

	return view
}

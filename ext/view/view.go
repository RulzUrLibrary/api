package view

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"io"
	"net/url"
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
	vars.Set("user", c.Get("user"))
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
		s, _ := a.Get(0).Interface().(string)
		return reflect.ValueOf(strings.Title(s))
	})
	view.AddGlobalFunc("capitalize", func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("capitalize", 1, 1)
		s, _ := a.Get(0).Interface().(string)
		if len(s) > 0 {
			s = strings.ToUpper(s[:1]) + s[1:]
		}
		return reflect.ValueOf(s)
	})
	view.AddGlobalFunc("query", func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("query", 3, 3)
		u, _ := a.Get(0).Interface().(*url.URL)
		v := u.Query()
		v.Set(a.Get(1).Interface().(string), fmt.Sprint(a.Get(2)))
		u.RawQuery = v.Encode()
		return reflect.ValueOf(u)
	})
	if config.Development {
		view.AddGlobalFunc("debug", func(a jet.Arguments) reflect.Value {
			a.RequireNumOfArguments("debug", 1, 1)
			return reflect.ValueOf(fmt.Sprintf("%#v", a.Get(0)))
		})
	}
	return view
}

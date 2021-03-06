package view

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/RulzUrLibrary/api/ext/i18n"
	"io"
	"net/url"
	"reflect"
	"strings"
)

type Configuration struct {
	Templates string
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
	vars.SetFunc("T", makeT(c))

	return tplt.Execute(w, vars, data)
}

func makeT(c echo.Context) jet.Func {
	t := i18n.GetI18n(c)
	return func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("t", 1, -1)
		id := a.Get(0).Interface().(string)
		args := []interface{}{}
		for i := 1; i < a.NumOfArguments(); i++ {
			args = append(args, a.Get(i))
		}
		return reflect.ValueOf(t(id, args...))
	}
}

func makeUrl(app *echo.Echo) jet.Func {
	return func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("url", 1, -1)
		name := a.Get(0).String()
		args := []interface{}{}
		for i := 1; i < a.NumOfArguments(); i++ {
			args = append(args, a.Get(i))
		}
		route := app.Reverse(name, args...)
		if route == "" {
			a.Panicf("route not found")
		}
		return reflect.ValueOf(route)
	}
}

func title(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("title", 1, 1)
	s, _ := a.Get(0).Interface().(string)
	return reflect.ValueOf(strings.Title(s))
}

func capitalize(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("capitalize", 1, 1)
	s, _ := a.Get(0).Interface().(string)
	if len(s) > 0 {
		s = strings.ToUpper(s[:1]) + s[1:]
	}
	return reflect.ValueOf(s)
}

func query(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("query", 3, 3)
	u, _ := a.Get(0).Interface().(*url.URL)
	v := u.Query()
	v.Set(a.Get(1).Interface().(string), fmt.Sprint(a.Get(2)))
	u.RawQuery = v.Encode()
	return reflect.ValueOf(u)
}

func str(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("str", 1, 1)
	return reflect.ValueOf(a.Get(0).Interface().(string))
}

func debug(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("debug", 1, 1)
	return reflect.ValueOf(fmt.Sprintf("%#v", a.Get(0)))
}

func New(app *echo.Echo, config Configuration) *View {
	view := &View{jet.NewHTMLSet(config.Templates)}
	view.SetDevelopmentMode(app.Debug)

	// add a reverse url helper
	view.AddGlobalFunc("url", makeUrl(app))
	view.AddGlobalFunc("title", title)
	view.AddGlobalFunc("capitalize", capitalize)
	view.AddGlobalFunc("query", query)
	view.AddGlobalFunc("str", str)

	if app.Debug {
		view.AddGlobalFunc("debug", debug)
	}
	return view
}

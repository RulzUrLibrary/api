package validator

import (
	"fmt"
	"github.com/labstack/echo"
	validate "gopkg.in/go-playground/validator.v9"
	"net/http"
	"reflect"
	"strings"
)

type Validator struct {
	validator *validate.Validate
}

var errors = map[string]string{
	"gt":  "%s value must be greater than %s, got '%d'",
	"lte": "%s value must be lower or equal to %s, got '%d'",
}

func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return err
	}
	payload := make(map[string][]string)
	ve, _ := err.(validate.ValidationErrors)
	for _, fe := range ve {
		msgs := []string{}

		if old, ok := payload[fe.Field()]; ok {
			msgs = old
		}
		msg := fmt.Sprintf(errors[fe.Tag()], fe.Field(), fe.Param(), fe.Value())
		payload[fe.Field()] = append(msgs, msg)
	}
	return echo.NewHTTPError(http.StatusBadRequest, payload)
}

func New() *Validator {
	validator := &Validator{validate.New()}
	validator.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return validator
}

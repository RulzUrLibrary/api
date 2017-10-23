package validator

import (
	"github.com/rulzurlibrary/api/utils"
	validate "gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
)

type Validator struct {
	validator *validate.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func Dump(err error, msgs map[string]utils.Dict) utils.Dict {
	ve, ok := err.(validate.ValidationErrors)
	if !ok {
		panic("dumping error failed")
	}
	payload := utils.Dict{}
	for _, fe := range ve {
		payload[fe.Field()] = msgs[fe.Field()][fe.Tag()]
	}
	return payload
}

func validatorGmail(fl validate.FieldLevel) bool {
	return utils.ValidMailProvider(fl.Field().String())
}

func New() *Validator {
	validator := &Validator{validate.New()}
	validator.validator.RegisterValidation("gmail", validatorGmail)
	validator.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})
	return validator
}

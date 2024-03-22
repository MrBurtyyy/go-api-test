package validation

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
	"strings"
)

func New() *validator.Validate {
	v := validator.New()
	err := v.RegisterValidation("dp", decimalPlaceValidator, false)
	if err != nil {
		panic("failed to add custom validation")
	}
	v.RegisterTagNameFunc(useJsonTag)
	return v
}

func decimalPlaceValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}

	switch field.Kind() {
	case reflect.Float64:
		s := strconv.FormatFloat(field.Float(), 'f', -1, 64)
		if i := strings.Index(s, "."); i > -1 {
			v := strings.Split(s, ".")
			if len(v) > 1 {
				return len(v[1]) <= param
			} else {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func useJsonTag(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

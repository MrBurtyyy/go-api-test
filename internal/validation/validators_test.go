package validation

import (
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"reflect"
	"testing"
)

type jsonTagTest struct {
	StructName string `json:"json_name"`
}

type decimalPlacesTest struct {
	NoDecimalPlaces float64 `validate:"dp=1"`
	Float64         float64 `validate:"dp=1"`
}

func TestDecimalPlaces(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("dp", decimalPlaceValidator)
	assert.Equal(t, err, nil)

	invalid := decimalPlacesTest{
		Float64: 1.123,
	}
	fieldsWithError := []string{
		"Float64",
	}

	errors := v.Struct(invalid).(validator.ValidationErrors)
	var fields []string
	for _, err := range errors {
		fields = append(fields, err.Field())
	}

	assert.Equal(t, fieldsWithError, fields)

	valid := decimalPlacesTest{
		NoDecimalPlaces: 1,
		Float64:         1.1,
	}

	err = v.Struct(valid)
	assert.Equal(t, err, nil)
}

func TestUseJsonName(t *testing.T) {
	typ := reflect.TypeOf(jsonTagTest{})
	field := typ.Field(0)

	jsonTag := useJsonTag(field)

	assert.Equal(t, jsonTag, "json_name")
}

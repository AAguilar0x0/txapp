package validatorv10

import (
	"github.com/AAguilar0x0/bapp/core/pkg/assert"
	"github.com/AAguilar0x0/bapp/core/services"
	"github.com/go-playground/validator/v10"
)

type ValidatorV10 struct {
	Validator *validator.Validate
}

func New() *ValidatorV10 {
	d := ValidatorV10{
		Validator: validator.New(),
	}
	err := d.Validator.RegisterValidation("enum_validation", func(fl validator.FieldLevel) bool {
		return fl.Field().Interface().(services.EnumValidator).ValidateEnum()
	})
	assert.NoError(err, "enum_validation registration", "faul", "RegisterValidation")
	return &d
}

func (d *ValidatorV10) Struct(s interface{}) error {
	return d.Validator.Struct(s)
}

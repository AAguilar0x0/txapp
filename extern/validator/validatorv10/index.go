package validatorv10

import (
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/go-playground/validator/v10"
)

type ValidatorV10 struct {
	init      bool
	Validator *validator.Validate
}

func (d *ValidatorV10) Init(env services.Environment) error {
	if d.init {
		return nil
	}
	d.Validator = validator.New()
	err := d.Validator.RegisterValidation("enum_validation", func(fl validator.FieldLevel) bool {
		return fl.Field().Interface().(services.EnumValidator).ValidateEnum()
	})
	d.init = true
	return err
}

func (*ValidatorV10) Close() {}

func (d *ValidatorV10) Struct(s interface{}) error {
	return d.Validator.Struct(s)
}

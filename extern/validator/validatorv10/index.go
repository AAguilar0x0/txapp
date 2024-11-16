package validatorv10

import (
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/go-playground/validator/v10"
)

type ValidatorV10 struct {
	Validator *validator.Validate
}

func New(env services.Environment) (*ValidatorV10, error) {
	d := ValidatorV10{validator.New()}
	err := d.Validator.RegisterValidation("enum_validation", func(fl validator.FieldLevel) bool {
		return fl.Field().Interface().(services.EnumValidator).ValidateEnum()
	})
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (*ValidatorV10) Close() error {
	return nil
}

func (d *ValidatorV10) Struct(s interface{}) *apierrors.APIError {
	err := d.Validator.Struct(s)
	if err != nil {
		return apierrors.BadRequest(err.Error())
	}
	return nil
}

func (d *ValidatorV10) Var(f interface{}, tag string) *apierrors.APIError {
	err := d.Validator.Var(f, tag)
	if err != nil {
		return apierrors.BadRequest(err.Error())
	}
	return nil
}

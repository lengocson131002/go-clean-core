package goplayaround

import (
	"github.com/go-playground/validator/v10"
	"github.com/lengocson131002/go-clean-core/validation"
)

type _validator struct {
	validator *validator.Validate
}

func NewGpValidator() validation.Validator {
	val := validator.New()
	return &_validator{
		validator: val,
	}
}

func (v *_validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err != nil {
		return err
	}
	return nil
}

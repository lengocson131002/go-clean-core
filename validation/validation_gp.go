package validation

import (
	"github.com/go-playground/validator/v10"
)

type _validator struct {
	validator *validator.Validate
}

func NewGpValidator() Validator {
	val := validator.New()
	return &_validator{
		validator: val,
	}
}

var _ Validator = (*_validator)(nil)

// Validate implements Validator.
func (v *_validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err != nil {
		return err
	}
	return nil
}

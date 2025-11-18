package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator instance and provides validation methods
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate validates a struct and returns a user-friendly error message
func (v *Validator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		firstError := validationErrors[0]
		errorMsg := fmt.Sprintf("Field '%s' failed validation '%s'",
			firstError.Field(),
			firstError.Tag())
		if firstError.Param() != "" {
			errorMsg += fmt.Sprintf(" (expected: %s)", firstError.Param())
		}
		return fmt.Errorf("%s", errorMsg)
	}
	return nil
}

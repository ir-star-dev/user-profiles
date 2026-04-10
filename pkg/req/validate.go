package req

import (
	"github.com/go-playground/validator/v10"
	"fmt"
)

type ValidationErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type ValidationErrors []ValidationErrorResponse

func (v ValidationErrors) Error() string {
	return "Validation failed"
}

func IsValid[T any](payload T) error {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		var errors ValidationErrors

		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationErrorResponse {
				Field: e.Field(),
				Tag:   buildMessage(e),
				Value: fmt.Sprintf("%v", e.Value()),
			})
		}

		return errors
	}

	return nil
}

func buildMessage(e validator.FieldError) string {
	switch e.Tag() {
		case "required":
			return "Required field"
		case "email":
			return "Incorrect email"
		case "min":
			return fmt.Sprintf("Min len %s", e.Param())
		case "max":
			return fmt.Sprintf("Max len %s", e.Param())
		default:
			return "Incorrect value"
	}
}
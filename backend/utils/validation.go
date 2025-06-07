package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []string {
	var errors []string
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s is required", strings.ToLower(e.Field())))
			case "min":
				errors = append(errors, fmt.Sprintf("%s must be at least %s characters", strings.ToLower(e.Field()), e.Param()))
			case "max":
				errors = append(errors, fmt.Sprintf("%s must be at most %s characters", strings.ToLower(e.Field()), e.Param()))
			case "gt":
				errors = append(errors, fmt.Sprintf("%s must be greater than %s", strings.ToLower(e.Field()), e.Param()))
			case "gte":
				errors = append(errors, fmt.Sprintf("%s must be greater than or equal to %s", strings.ToLower(e.Field()), e.Param()))
			case "lt":
				errors = append(errors, fmt.Sprintf("%s must be less than %s", strings.ToLower(e.Field()), e.Param()))
			case "lte":
				errors = append(errors, fmt.Sprintf("%s must be less than or equal to %s", strings.ToLower(e.Field()), e.Param()))
			case "email":
				errors = append(errors, fmt.Sprintf("%s must be a valid email address", strings.ToLower(e.Field())))
			case "hexcolor":
				errors = append(errors, fmt.Sprintf("%s must be a valid hex color", strings.ToLower(e.Field())))
			case "uuid":
				errors = append(errors, fmt.Sprintf("%s must be a valid UUID", strings.ToLower(e.Field())))
			default:
				errors = append(errors, fmt.Sprintf("%s is invalid", strings.ToLower(e.Field())))
			}
		}
	} else {
		errors = append(errors, err.Error())
	}
	
	return errors
}

package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string {
	return "validation failed"
}

func New() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]; name != "" && name != "-" {
			return name
		}
		if name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]; name != "" && name != "-" {
			return name
		}
		return fld.Name
	})
	return &Validator{v: v}
}

func (v *Validator) Struct(payload any) error {
	if err := v.v.Struct(payload); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			fields := make(map[string]string, len(verr))
			for _, fe := range verr {
				field := fe.Field()
				fields[field] = messageForTag(fe)
			}
			return ValidationError{Fields: fields}
		}
		return err
	}
	return nil
}

func messageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	default:
		return fmt.Sprintf("failed validation %s", fe.Tag())
	}
}

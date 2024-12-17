package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Validator struct{}

// NewValidator creates a new Validator instance.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateStruct validates a struct based on annotation tags.
func (v *Validator) ValidateStruct(s interface{}) error {
	val := reflect.ValueOf(s)
	typ := val.Type()

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("validation only works on structs; got %T", s)
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			if err := v.validateField(field.String(), rule); err != nil {
				return fmt.Errorf("field %s: %w", fieldType.Name, err)
			}
		}
	}
	return nil
}

// validateField applies validation rules to a single field
func (v *Validator) validateField(value, rule string) error {
	parts := strings.Split(rule, "=")
	ruleName := parts[0]

	switch ruleName {
	case "required":
		if value == "" {
			return fmt.Errorf("field is required")
		}
	case "email":
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			return fmt.Errorf("invalid email format")
		}
	case "min":
		if len(parts) != 2 {
			return fmt.Errorf("min rule requires a length parameter")
		}
		minLen, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid length")
		}
		if len(value) < minLen {
			return fmt.Errorf("minimum length is %d", minLen)
		}
	}
	return nil
}

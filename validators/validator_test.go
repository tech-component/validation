package validators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		input   interface{}
		wantErr error
	}{
		{
			name: "valid struct",
			input: struct {
				RequiredField string `validate:"required"`
				EmailField    string `validate:"email"`
				PasswordField string `validate:"min=8"`
				NoValidation  string
			}{
				RequiredField: "present",
				EmailField:    "test@example.com",
				PasswordField: "password123",
			},
			wantErr: nil,
		},
		{
			name: "missing required field",
			input: struct {
				RequiredField string `validate:"required"`
			}{},
			wantErr: fmt.Errorf("field RequiredField: field is required"),
		},
		{
			name: "invalid email",
			input: struct {
				EmailField string `validate:"email"`
			}{
				EmailField: "invalid-email",
			},
			wantErr: fmt.Errorf("field EmailField: invalid email format"),
		},
		{
			name: "password too short",
			input: struct {
				PasswordField string `validate:"min=8"`
			}{
				PasswordField: "short",
			},
			wantErr: fmt.Errorf("field PasswordField: minimum length is 8"),
		},
		{
			name:    "non-struct input",
			input:   "not a struct",
			wantErr: fmt.Errorf("validation only works on structs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStruct(tt.input)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateField(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name       string
		value      string
		rule       string
		wantErrStr string
	}{
		{
			name:  "required - valid",
			value: "present",
			rule:  "required",
		},
		{
			name:       "required - invalid",
			value:      "",
			rule:       "required",
			wantErrStr: "field is required",
		},
		{
			name:  "email - valid",
			value: "test@example.com",
			rule:  "email",
		},
		{
			name:       "email - invalid",
			value:      "not-an-email",
			rule:       "email",
			wantErrStr: "invalid email format",
		},
		{
			name:  "min - valid",
			value: "password123",
			rule:  "min=8",
		},
		{
			name:       "min - invalid",
			value:      "short",
			rule:       "min=8",
			wantErrStr: "minimum length is 8",
		},
		{
			name:       "min - missing parameter",
			value:      "any",
			rule:       "min",
			wantErrStr: "min rule requires a length parameter",
		},
		{
			name:       "min - invalid parameter",
			value:      "any",
			rule:       "min=invalid",
			wantErrStr: "invalid length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateField(tt.value, tt.rule)
			if tt.wantErrStr != "" {
				assert.ErrorContains(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

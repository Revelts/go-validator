package validator

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field     string
	FieldPath string
	Message   string
}

func (e *ValidationError) Error() string {
	if e.FieldPath != "" && e.FieldPath != e.Field {
		return fmt.Sprintf("%s: %s", e.FieldPath, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []*ValidationError
}

func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d validation errors:\n", len(e.Errors)))
	for i, err := range e.Errors {
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}
	return b.String()
}

// Add adds a validation error
func (e *ValidationErrors) Add(field, fieldPath, message string) {
	e.Errors = append(e.Errors, &ValidationError{
		Field:     field,
		FieldPath: fieldPath,
		Message:   message,
	})
}

// HasErrors returns true if there are any errors
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// Error constants
const (
	ErrInvalidStruct             = "not a struct type or memory address of struct"
	ErrRequiredString            = "value is required"
	ErrRequiredNumber            = "value is required, minimum 1"
	ErrRequiredSlice             = "value is required, minimum 1 value"
	ErrInvalidNumberType         = "invalid type of number"
	ErrInvalidTextType           = "invalid type of text"
	ErrInvalidEmailFormat        = "value format is not email"
	ErrInvalidAlphabetFormat     = "value format is not alphabet"
	ErrInvalidAlphanumericFormat = "value format is not alphanumeric"
	ErrMinValueString            = "should be greater than %s %s"
	ErrMinValueNumber            = "should be greater than %s"
	ErrMaxValueString            = "should be lower than %s %s"
	ErrMaxValueNumber            = "should be lower than %s"
	ErrInvalidRegexFormat        = "invalid value of expression"
	ErrStartsWith                = "value should be starts with %s"
	ErrEndsWith                  = "value should be ends with %s"
	ErrValueOf                   = "value should be contains one of %s"
)

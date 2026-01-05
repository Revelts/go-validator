package validator

import (
	"fmt"
	"reflect"
	"strings"
)

// validateMax validates maximum value/length
func validateMax(ctx *ValidationContext) error {
	tagVal, err := getInt(ctx.TagValue)
	if err != nil {
		return fmt.Errorf("invalid max value: %w", err)
	}

	// Use reflect.Value directly to avoid type assertions where possible
	rv := ctx.ReflectValue

	switch ctx.ValueType {
	case reflect.String:
		if rv.Len() > tagVal {
			return fmt.Errorf(ErrMaxValueString, ctx.TagValue, "char")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Int() > int64(tagVal) {
			return fmt.Errorf(ErrMaxValueNumber, ctx.TagValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rv.Uint() > uint64(tagVal) {
			return fmt.Errorf(ErrMaxValueNumber, ctx.TagValue)
		}
	case reflect.Float32, reflect.Float64:
		if rv.Float() > float64(tagVal) {
			return fmt.Errorf(ErrMaxValueNumber, ctx.TagValue)
		}
	}

	return nil
}

// validateMin validates minimum value/length
func validateMin(ctx *ValidationContext) error {
	tagVal, err := getInt(ctx.TagValue)
	if err != nil {
		return fmt.Errorf("invalid min value: %w", err)
	}

	// Use reflect.Value directly to avoid type assertions where possible
	rv := ctx.ReflectValue

	switch ctx.ValueType {
	case reflect.String:
		if rv.Len() < tagVal {
			return fmt.Errorf(ErrMinValueString, ctx.TagValue, "char")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Int() < int64(tagVal) {
			return fmt.Errorf(ErrMinValueNumber, ctx.TagValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rv.Uint() < uint64(tagVal) {
			return fmt.Errorf(ErrMinValueNumber, ctx.TagValue)
		}
	case reflect.Float32, reflect.Float64:
		if rv.Float() < float64(tagVal) {
			return fmt.Errorf(ErrMinValueNumber, ctx.TagValue)
		}
	}

	return nil
}

// validateField validates required fields
func validateField(ctx *ValidationContext) error {
	if ctx.TagValue != tagRequired {
		return nil
	}

	rv := ctx.ReflectValue

	switch ctx.ValueType {
	case reflect.String:
		val := strings.TrimSpace(rv.String())
		if val == "" {
			return fmt.Errorf(ErrRequiredString)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Int() == 0 {
			return fmt.Errorf(ErrRequiredNumber)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rv.Uint() == 0 {
			return fmt.Errorf(ErrRequiredNumber)
		}
	case reflect.Float32, reflect.Float64:
		if rv.Float() == 0 {
			return fmt.Errorf(ErrRequiredNumber)
		}
	case reflect.Slice, reflect.Array:
		if rv.Len() == 0 {
			return fmt.Errorf(ErrRequiredSlice)
		}
	case reflect.Map:
		if rv.Len() == 0 {
			return fmt.Errorf(ErrRequiredSlice)
		}
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return fmt.Errorf(ErrRequiredString)
		}
	}

	return nil
}

// validateStartsWith validates string prefix
func validateStartsWith(ctx *ValidationContext) error {
	if ctx.ValueType != reflect.String {
		return fmt.Errorf(ErrInvalidTextType)
	}

	val := ctx.Value.(string)
	if val != "" && !strings.HasPrefix(val, ctx.TagValue) {
		return fmt.Errorf(ErrStartsWith, ctx.TagValue)
	}

	return nil
}

// validateEndsWith validates string suffix
func validateEndsWith(ctx *ValidationContext) error {
	if ctx.ValueType != reflect.String {
		return fmt.Errorf(ErrInvalidTextType)
	}

	val := ctx.Value.(string)
	if val != "" && !strings.HasSuffix(val, ctx.TagValue) {
		return fmt.Errorf(ErrEndsWith, ctx.TagValue)
	}

	return nil
}

// validateValueOf validates that value is one of the allowed values
func validateValueOf(ctx *ValidationContext) error {
	if ctx.ValueType != reflect.String {
		return fmt.Errorf(ErrInvalidTextType)
	}

	val := ctx.Value.(string)
	if val == "" {
		return nil
	}

	allowedValues := strings.Split(ctx.TagValue, ",")
	for _, allowed := range allowedValues {
		if val == strings.TrimSpace(allowed) {
			return nil
		}
	}

	return fmt.Errorf(ErrValueOf, ctx.TagValue)
}

// validateFormat validates string format (email, alphabet, alphanumeric)
func validateFormat(ctx *ValidationContext) error {
	if ctx.ValueType != reflect.String {
		return fmt.Errorf(ErrInvalidTextType)
	}

	val := ctx.Value.(string)
	if val == "" {
		return nil
	}

	var pattern string
	var errMsg string

	switch ctx.TagValue {
	case tagEmail:
		pattern = mailFormatRegex
		errMsg = ErrInvalidEmailFormat
	case tagAlphanumeric:
		pattern = alphanumericFormatRegex
		errMsg = ErrInvalidAlphanumericFormat
	case tagAlphabet:
		pattern = alphabetFormatRegex
		errMsg = ErrInvalidAlphabetFormat
	default:
		return fmt.Errorf("%s: unknown format", ctx.TagValue)
	}

	re, err := getRegex(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	if !re.MatchString(val) {
		return &formatError{message: errMsg}
	}

	return nil
}

// validateFieldType validates field type
func validateFieldType(ctx *ValidationContext) error {
	switch ctx.TagValue {
	case tagText:
		if ctx.ValueType != reflect.String {
			return fmt.Errorf(ErrInvalidTextType)
		}
	case tagNumber:
		if ctx.ValueType == reflect.String {
			re, err := getRegex(numericFormatRegex)
			if err != nil {
				return fmt.Errorf("invalid regex pattern: %w", err)
			}
			if !re.MatchString(ctx.Value.(string)) {
				return fmt.Errorf(ErrInvalidNumberType)
			}
			return nil
		}

		switch ctx.ValueType {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			// Valid number types
		default:
			return fmt.Errorf(ErrInvalidNumberType)
		}
	}

	return nil
}

// validateMatch validates against a regex pattern
func validateMatch(ctx *ValidationContext) error {
	if ctx.ValueType != reflect.String {
		return fmt.Errorf(ErrInvalidTextType)
	}

	val := ctx.Value.(string)
	if val == "" {
		return nil
	}

	re, err := getRegex(ctx.TagValue)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	if !re.MatchString(val) {
		return fmt.Errorf(ErrInvalidRegexFormat)
	}

	return nil
}


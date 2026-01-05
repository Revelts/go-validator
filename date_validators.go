package validator

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Common date formats to try when parsing
var dateFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006/01/02",
	"2006/01/02 15:04:05",
	"01/02/2006",
	"01-02-2006",
	"02/01/2006",
	"02-01-2006",
}

// parseDateString attempts to parse a date string using common formats
func parseDateString(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	var lastErr error
	for _, format := range dateFormats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	return time.Time{}, fmt.Errorf("unable to parse date '%s': %w", dateStr, lastErr)
}

// getFieldValueFromStruct gets a field value from a struct by field name
func getFieldValueFromStruct(structValue reflect.Value, fieldName string) (reflect.Value, bool) {
	structType := structValue.Type()
	
	// Try exact match first
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Name == fieldName {
			if structValue.Field(i).CanInterface() {
				return structValue.Field(i), true
			}
		}
		// Also check JSON tag
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			jsonName := jsonTag
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				jsonName = jsonTag[:idx]
			}
			if jsonName == fieldName {
				if structValue.Field(i).CanInterface() {
					return structValue.Field(i), true
				}
			}
		}
	}
	
	return reflect.Value{}, false
}

// validateStartDate validates that start date is less than end date
// Tag value should be the name of the end date field to compare against
func validateStartDate(ctx *ValidationContext) error {
	if ctx.ValueType != reflect.String {
		return fmt.Errorf(ErrInvalidTextType)
	}

	startDateStr := ctx.Value.(string)
	if startDateStr == "" {
		return nil // Empty values are handled by required validator
	}

	// Parse start date
	startDate, err := parseDateString(startDateStr)
	if err != nil {
		return fmt.Errorf("invalid start date format: %w", err)
	}

	// Get the end date field name from tag value
	endDateFieldName := ctx.TagValue
	if endDateFieldName == "" {
		return fmt.Errorf("startdate tag must specify the end date field name")
	}

	// Get parent struct from context
	parentStructValue := getParentStructValue(ctx)
	if !parentStructValue.IsValid() {
		return fmt.Errorf("unable to access parent struct for field comparison")
	}

	// Get end date field value
	endDateValue, found := getFieldValueFromStruct(parentStructValue, endDateFieldName)
	if !found {
		return fmt.Errorf("end date field '%s' not found in struct", endDateFieldName)
	}

	// Check if end date is empty (skip validation if empty)
	if endDateValue.Kind() == reflect.String {
		endDateStr := endDateValue.String()
		if strings.TrimSpace(endDateStr) == "" {
			return nil // End date is empty, skip comparison
		}

		// Parse end date
		endDate, err := parseDateString(endDateStr)
		if err != nil {
			return nil // If end date is invalid, let enddate validator handle it
		}

		// Validate: StartDate must be less than EndDate (start comes before end)
		// Note: If you need StartDate > EndDate, reverse this logic
		if startDate.After(endDate) {
			return fmt.Errorf("start date must be before or equal to end date")
		}
	}

	return nil
}

// validateEndDate validates that end date is greater than start date
// Tag value should be the name of the start date field to compare against
func validateEndDate(ctx *ValidationContext) error {
	if ctx.ValueType != reflect.String {
		return fmt.Errorf(ErrInvalidTextType)
	}

	endDateStr := ctx.Value.(string)
	if endDateStr == "" {
		return nil // Empty values are handled by required validator
	}

	// Parse end date
	endDate, err := parseDateString(endDateStr)
	if err != nil {
		return fmt.Errorf("invalid end date format: %w", err)
	}

	// Get the start date field name from tag value
	startDateFieldName := ctx.TagValue
	if startDateFieldName == "" {
		return fmt.Errorf("enddate tag must specify the start date field name")
	}

	// Get parent struct from context
	parentStructValue := getParentStructValue(ctx)
	if !parentStructValue.IsValid() {
		return fmt.Errorf("unable to access parent struct for field comparison")
	}

	// Get start date field value
	startDateValue, found := getFieldValueFromStruct(parentStructValue, startDateFieldName)
	if !found {
		return fmt.Errorf("start date field '%s' not found in struct", startDateFieldName)
	}

	// Check if start date is empty (skip validation if empty)
	if startDateValue.Kind() == reflect.String {
		startDateStr := startDateValue.String()
		if strings.TrimSpace(startDateStr) == "" {
			return nil // Start date is empty, skip comparison
		}

		// Parse start date
		startDate, err := parseDateString(startDateStr)
		if err != nil {
			return nil // If start date is invalid, let startdate validator handle it
		}

		// Validate: EndDate must be greater than StartDate (end comes after start)
		if endDate.Before(startDate) {
			return fmt.Errorf("end date must be after or equal to start date")
		}
	}

	return nil
}

// getParentStructValue gets the parent struct value from context
func getParentStructValue(ctx *ValidationContext) reflect.Value {
	if ctx.ParentStruct.IsValid() {
		return ctx.ParentStruct
	}
	return reflect.Value{}
}


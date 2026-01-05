package validator

import (
	"reflect"
	"strings"
)

// ValidationOptions configures validation behavior
type ValidationOptions struct {
	CollectAllErrors bool   // Collect all errors instead of returning first
	FieldPathPrefix  string // Prefix for field paths in errors
}

// DefaultOptions returns default validation options
func DefaultOptions() *ValidationOptions {
	return &ValidationOptions{
		CollectAllErrors: true,
		FieldPathPrefix:  "",
	}
}

// Validate validates a struct or slice of structs
func Validate(data interface{}) error {
	return ValidateWithOptions(data, DefaultOptions())
}

// ValidateWithOptions validates with custom options
func ValidateWithOptions(data interface{}, opts *ValidationOptions) error {
	if opts == nil {
		opts = DefaultOptions()
	}

	registry := GetRegistry()
	errors := &ValidationErrors{Errors: make([]*ValidationError, 0)}

	field := reflect.TypeOf(data)
	value := reflect.ValueOf(data)

	if field.Kind() != reflect.Struct && field.Kind() != reflect.Slice {
		return &ValidationError{
			Field:   "root",
			Message: ErrInvalidStruct,
		}
	}

	if field.Kind() == reflect.Slice {
		if err := validateSlice(value, registry, opts, errors, ""); err != nil && !opts.CollectAllErrors {
			return err
		}
		if errors.HasErrors() {
			return errors
		}
		return nil
	}

	validateStruct(field, value, registry, opts, errors, opts.FieldPathPrefix)

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// validateStruct validates a struct recursively
func validateStruct(
	field reflect.Type,
	value reflect.Value,
	registry *ValidatorRegistry,
	opts *ValidationOptions,
	errors *ValidationErrors,
	fieldPath string,
) {
	for i := 0; i < field.NumField(); i++ {
		dataField := field.Field(i)
		fieldName := getFieldName(dataField)
		currentPath := buildFieldPath(fieldPath, fieldName)

		if !value.Field(i).CanInterface() {
			continue
		}

		dataValue := value.Field(i)
		kind := dataField.Type.Kind()

		switch kind {
		case reflect.Map, reflect.Slice, reflect.Array:
			if err := validateMapSlice(dataField, dataValue, value, field, registry, opts, errors, currentPath); err != nil && !opts.CollectAllErrors {
				errors.Add(fieldName, currentPath, err.Error())
				if !opts.CollectAllErrors {
					return
				}
			}
			continue

		case reflect.Struct:
			if dataValue.CanInterface() {
				validateStruct(dataField.Type, dataValue, registry, opts, errors, currentPath)
			}
			continue

		case reflect.Ptr:
			if dataValue.IsNil() {
				// Validate required tag even for nil pointers
				if err := validateFieldTags(dataField, reflect.ValueOf(nil), value, field, registry, errors, fieldName, currentPath); err != nil && !opts.CollectAllErrors {
					return
				}
				continue
			}

			elemValue := dataValue.Elem()
			if elemValue.Kind() == reflect.Struct {
				validateStruct(elemValue.Type(), elemValue, registry, opts, errors, currentPath)
			} else {
				// Validate pointer to primitive
				if err := validateFieldTags(dataField, elemValue, value, field, registry, errors, fieldName, currentPath); err != nil && !opts.CollectAllErrors {
					return
				}
			}
			continue
		}

		if err := validateFieldTags(dataField, dataValue, value, field, registry, errors, fieldName, currentPath); err != nil && !opts.CollectAllErrors {
			return
		}
	}
}

// Built-in tag names - defined once to avoid allocations
var builtinTagNames = []string{"max", "min", "field", "startswith", "endswith", "value_of", "format", "fieldType", "match", "startdate", "enddate"}

// validateFieldTags validates all tags on a field
func validateFieldTags(
	dataField reflect.StructField,
	dataValue reflect.Value,
	parentStruct reflect.Value,
	parentStructType reflect.Type,
	registry *ValidatorRegistry,
	errors *ValidationErrors,
	fieldName, fieldPath string,
) error {
	// Cache interface conversion to avoid multiple calls
	var cachedValue interface{}
	var valueCached bool

	for _, tagName := range builtinTagNames {
		tagValue, ok := dataField.Tag.Lookup(tagName)
		if !ok || tagValue == "" {
			continue
		}

		validator, exists := registry.Get(tagName)
		if !exists {
			continue
		}

		// Lazy load interface value only when needed
		if !valueCached && dataValue.CanInterface() {
			cachedValue = dataValue.Interface()
			valueCached = true
		}

		ctx := &ValidationContext{
			FieldName:        fieldName,
			FieldPath:        fieldPath,
			TagValue:         tagValue,
			Value:            cachedValue,
			ValueType:        dataValue.Kind(),
			ReflectValue:     dataValue,
			StructField:      dataField,
			ParentStruct:     parentStruct,
			ParentStructType: parentStructType,
		}

		if err := validator(ctx); err != nil {
			errors.Add(fieldName, fieldPath, err.Error())
			return err
		}
	}

	return nil
}

// validateMapSlice validates map, slice, or array fields
func validateMapSlice(
	dataField reflect.StructField,
	value reflect.Value,
	parentStruct reflect.Value,
	parentStructType reflect.Type,
	registry *ValidatorRegistry,
	opts *ValidationOptions,
	errors *ValidationErrors,
	fieldPath string,
) error {
	// Validate required tag on the collection itself
	if err := validateFieldTags(dataField, value, parentStruct, parentStructType, registry, errors, dataField.Name, fieldPath); err != nil && !opts.CollectAllErrors {
		return err
	}

	switch value.Kind() {
	case reflect.Map:
		for _, key := range value.MapKeys() {
			elemValue := value.MapIndex(key)
			if elemValue.CanInterface() {
				if err := validateFieldTags(dataField, elemValue, parentStruct, parentStructType, registry, errors, dataField.Name, fieldPath); err != nil && !opts.CollectAllErrors {
					return err
				}
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			elemValue := value.Index(i)
			if !elemValue.CanInterface() {
				continue
			}

			if elemValue.Kind() == reflect.Struct {
				validateStruct(elemValue.Type(), elemValue, registry, opts, errors, fieldPath)
			} else {
				if err := validateFieldTags(dataField, elemValue, parentStruct, parentStructType, registry, errors, dataField.Name, fieldPath); err != nil && !opts.CollectAllErrors {
					return err
				}
			}
		}
	}

	return nil
}

// validateSlice validates a slice of structs
func validateSlice(
	value reflect.Value,
	registry *ValidatorRegistry,
	opts *ValidationOptions,
	errors *ValidationErrors,
	fieldPath string,
) error {
	for i := 0; i < value.Len(); i++ {
		elemValue := value.Index(i)
		if !elemValue.CanInterface() {
			continue
		}

		if elemValue.Kind() == reflect.Struct {
			validateStruct(elemValue.Type(), elemValue, registry, opts, errors, fieldPath)
		} else if elemValue.Kind() == reflect.Ptr && !elemValue.IsNil() {
			elem := elemValue.Elem()
			if elem.Kind() == reflect.Struct {
				validateStruct(elem.Type(), elem, registry, opts, errors, fieldPath)
			}
		}
	}

	return nil
}

// getFieldName extracts field name, preferring JSON tag
func getFieldName(field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		// Handle json tag with options like "name,omitempty"
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			return jsonTag[:idx]
		}
		return jsonTag
	}
	return field.Name
}

// buildFieldPath builds a field path for nested structures
func buildFieldPath(parent, current string) string {
	if parent == "" {
		return current
	}
	return parent + "." + current
}

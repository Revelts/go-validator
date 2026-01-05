package validator

import (
	"reflect"
	"sync"
)

// ValidatorFunc is a function type for custom validators
type ValidatorFunc func(ctx *ValidationContext) error

// ValidationContext provides context for validation
type ValidationContext struct {
	FieldName      string
	FieldPath      string
	TagValue       string
	Value          interface{}
	ValueType      reflect.Kind
	ReflectValue   reflect.Value
	StructField    reflect.StructField
	ParentStruct   reflect.Value // Parent struct value for accessing sibling fields
	ParentStructType reflect.Type // Parent struct type
}

// ValidatorRegistry manages all validators
type ValidatorRegistry struct {
	validators map[string]ValidatorFunc
	mu         sync.RWMutex
}

var defaultRegistry *ValidatorRegistry
var registryOnce sync.Once

// GetRegistry returns the default validator registry
func GetRegistry() *ValidatorRegistry {
	registryOnce.Do(func() {
		defaultRegistry = NewRegistry()
		defaultRegistry.registerBuiltinValidators()
	})
	return defaultRegistry
}

// NewRegistry creates a new validator registry
func NewRegistry() *ValidatorRegistry {
	return &ValidatorRegistry{
		validators: make(map[string]ValidatorFunc),
	}
}

// Register adds a custom validator to the registry
func (r *ValidatorRegistry) Register(name string, fn ValidatorFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.validators[name] = fn
}

// Get retrieves a validator by name
func (r *ValidatorRegistry) Get(name string) (ValidatorFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, ok := r.validators[name]
	return fn, ok
}

// registerBuiltinValidators registers all built-in validators
func (r *ValidatorRegistry) registerBuiltinValidators() {
	r.Register("max", validateMax)
	r.Register("min", validateMin)
	r.Register("field", validateField)
	r.Register("startswith", validateStartsWith)
	r.Register("endswith", validateEndsWith)
	r.Register("value_of", validateValueOf)
	r.Register("format", validateFormat)
	r.Register("fieldType", validateFieldType)
	r.Register("match", validateMatch)
	r.Register("startdate", validateStartDate)
	r.Register("enddate", validateEndDate)
}

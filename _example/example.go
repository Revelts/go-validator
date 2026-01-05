package main

import (
	"fmt"
	"validator"
)

// ============================================================================
// COMPREHENSIVE EXAMPLES - All Available Validators
// ============================================================================

// Example 1: Basic Validations
type BasicExample struct {
	// field: required validation
	RequiredString string   `field:"required"`
	RequiredNumber int      `field:"required"`
	RequiredSlice  []string `field:"required"`

	// min/max: length and value constraints
	MinLengthString string  `min:"5"`              // Minimum 5 characters
	MaxLengthString string  `max:"10"`             // Maximum 10 characters
	MinMaxString    string  `min:"3" max:"20"`     // Between 3-20 characters
	MinValue        int     `min:"1"`              // Minimum value 1
	MaxValue        int     `max:"100"`            // Maximum value 100
	MinMaxValue     float64 `min:"0.5" max:"99.9"` // Between 0.5-99.9

	// fieldType: type validation
	TextField   string `fieldType:"text"`   // Must be text
	NumberField string `fieldType:"number"` // Must be numeric string
	IntField    int    `fieldType:"number"` // Must be number type
}

// Example 2: String Format Validations
type FormatExample struct {
	// format: email, alphabet, alphanumeric
	Email        string `format:"email" field:"required"`
	AlphabetOnly string `format:"alphabet"`     // Only letters and spaces
	Alphanumeric string `format:"alphanumeric"` // Letters, numbers, and spaces

	// startswith/endswith: prefix and suffix validation
	StartsWithPrefix string `startswith:"https://"`             // Must start with https://
	EndsWithSuffix   string `endswith:".com"`                   // Must end with .com
	BothPrefixSuffix string `startswith:"user_" endswith:"_id"` // Both conditions
}

// Example 3: Value Constraints
type ValueExample struct {
	// value_of: must be one of the specified values
	Status      string `value_of:"active,inactive,pending" field:"required"`
	Priority    string `value_of:"low,medium,high"`
	CountryCode string `value_of:"US,UK,CA,AU"`

	// match: custom regex pattern
	PhoneNumber   string `match:"^[0-9]{10}$"`            // Exactly 10 digits
	ZipCode       string `match:"^[0-9]{5}(-[0-9]{4})?$"` // US zip code format
	CustomPattern string `match:"^[A-Z]{2}[0-9]{4}$"`     // 2 letters + 4 digits
}

// Example 4: Date Validations
type DateExample struct {
	// startdate/enddate: date comparison
	StartDate string `startdate:"EndDate" field:"required"` // Must be before EndDate
	EndDate   string `enddate:"StartDate" field:"required"` // Must be after StartDate
	EventName string `field:"required"`
}

// Example 5: Combined Validations
type CombinedExample struct {
	// Multiple validations on single field
	Username string `field:"required" min:"3" max:"20" format:"alphanumeric"`
	Email    string `field:"required" format:"email" max:"100"`
	Age      int    `field:"required" min:"18" max:"120" fieldType:"number"`
	Password string `field:"required" min:"8" max:"50" match:"^[A-Za-z0-9@#$%^&*!]+$"`
}

// Example 6: Nested Structures
type Address struct {
	Street  string `field:"required" format:"alphanumeric"`
	City    string `field:"required" format:"alphabet" min:"2" max:"50"`
	ZipCode string `field:"required" match:"^[0-9]{5}$"`
	Country string `value_of:"US,UK,CA,AU,MX" field:"required"`
}

type Person struct {
	Name    string  `field:"required" min:"2" max:"50" format:"alphabet"`
	Email   string  `field:"required" format:"email"`
	Age     int     `field:"required" min:"18" max:"120"`
	Address Address // Nested struct validation
}

// Example 7: Pointers and Optional Fields
type OptionalExample struct {
	RequiredField  string  `field:"required"`
	OptionalString *string `min:"5" max:"20"`  // Optional, but if provided must be 5-20 chars
	OptionalNumber *int    `min:"1" max:"100"` // Optional, but if provided must be 1-100
	OptionalEmail  *string `format:"email"`    // Optional email
}

// Example 8: Slices and Arrays
type SliceExample struct {
	RequiredSlice []string `field:"required"` // Must have at least 1 item
	Tags          []string `min:"1" max:"10"` // Each tag 1-10 chars
	Numbers       []int    `field:"required"` // Required slice
	OptionalSlice []string // Optional slice
}

// Example 9: Maps
type MapExample struct {
	RequiredMap map[string]string `field:"required"`   // Map must not be empty
	ConfigMap   map[string]string `fieldType:"number"` // Values must be numeric
}

// Example 10: Complex Real-World Example
type UserRegistration struct {
	// Basic Info
	Username    string `field:"required" min:"3" max:"20" format:"alphanumeric"`
	Email       string `field:"required" format:"email" max:"100"`
	Password    string `field:"required" min:"8" max:"50" match:"^[A-Za-z0-9@#$%^&*!]+$"`
	ConfirmPass string `field:"required"` // Note: password match would need custom validator

	// Personal Info
	FirstName string `field:"required" min:"2" max:"50" format:"alphabet"`
	LastName  string `field:"required" min:"2" max:"50" format:"alphabet"`
	Age       int    `field:"required" min:"18" max:"120" fieldType:"number"`

	// Contact
	Phone   string `field:"required" match:"^[0-9]{10}$" fieldType:"number"`
	Website string `startswith:"https://"` // Optional, but if provided must start with https://

	// Preferences
	Status  string `value_of:"active,inactive" field:"required"`
	Country string `value_of:"US,UK,CA,AU,MX,DE,FR" field:"required"`

	// Dates
	StartDate string `startdate:"EndDate" field:"required"`
	EndDate   string `enddate:"StartDate" field:"required"`

	// Address
	Address Address

	// Optional Fields
	Bio  *string  `max:"500"`        // Optional bio, max 500 chars
	Tags []string `min:"1" max:"20"` // Each tag 1-20 chars
}

// ============================================================================
// DEMONSTRATION FUNCTIONS
// ============================================================================

func main() {
	fmt.Println("=== GoValidator - Comprehensive Examples ===\n")

	// Example 1: Basic Validations
	demonstrateBasicValidations()

	// Example 2: Format Validations
	demonstrateFormatValidations()

	// Example 3: Value Constraints
	demonstrateValueConstraints()

	// Example 4: Date Validations
	demonstrateDateValidations()

	// Example 5: Combined Validations
	demonstrateCombinedValidations()

	// Example 6: Nested Structures
	demonstrateNestedStructures()

	// Example 7: Optional Fields (Pointers)
	demonstrateOptionalFields()

	// Example 8: Slices
	demonstrateSlices()

	// Example 9: Maps
	demonstrateMaps()

	// Example 10: Real-World Example
	demonstrateRealWorldExample()
}

func demonstrateBasicValidations() {
	fmt.Println("1. BASIC VALIDATIONS")
	fmt.Println("-------------------")

	// Valid example
	valid := BasicExample{
		RequiredString:  "Hello",
		RequiredNumber:  42,
		RequiredSlice:   []string{"item1", "item2"},
		MinLengthString: "Hello World", // 11 chars, >= 5 ✓
		MaxLengthString: "Short",       // 5 chars, <= 10 ✓
		MinMaxString:    "Valid",       // 5 chars, 3-20 ✓
		MinValue:        10,            // >= 1 ✓
		MaxValue:        50,            // <= 100 ✓
		MinMaxValue:     50.5,          // 0.5-99.9 ✓
		TextField:       "text",
		NumberField:     "12345",
		IntField:        100,
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid basic example passed")
	}

	// Invalid example - missing required field
	invalid := BasicExample{
		RequiredString: "", // Missing required field
		RequiredNumber: 0,  // Missing required field
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught validation error (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateFormatValidations() {
	fmt.Println("2. FORMAT VALIDATIONS")
	fmt.Println("---------------------")

	valid := FormatExample{
		Email:            "user@example.com",
		AlphabetOnly:     "Hello World",
		Alphanumeric:     "Hello123 World",
		StartsWithPrefix: "https://example.com",
		EndsWithSuffix:   "example.com",
		BothPrefixSuffix: "user_123_id",
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid format example passed")
	}

	invalid := FormatExample{
		Email: "invalid-email", // Invalid email format
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught format error (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateValueConstraints() {
	fmt.Println("3. VALUE CONSTRAINTS")
	fmt.Println("--------------------")

	valid := ValueExample{
		Status:        "active",
		Priority:      "high",
		CountryCode:   "US",
		PhoneNumber:   "1234567890",
		ZipCode:       "12345-6789",
		CustomPattern: "AB1234",
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid value constraints example passed")
	}

	invalid := ValueExample{
		Status: "invalid_status", // Not in allowed values
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught value constraint error (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateDateValidations() {
	fmt.Println("4. DATE VALIDATIONS")
	fmt.Println("-------------------")

	valid := DateExample{
		StartDate: "2024-01-01",
		EndDate:   "2024-01-31",
		EventName: "Conference",
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid date example passed (StartDate < EndDate)")
	}

	invalid := DateExample{
		StartDate: "2024-01-31", // After end date
		EndDate:   "2024-01-01",
		EventName: "Conference",
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught date validation error (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateCombinedValidations() {
	fmt.Println("5. COMBINED VALIDATIONS")
	fmt.Println("-----------------------")

	valid := CombinedExample{
		Username: "john_doe123",
		Email:    "john@example.com",
		Age:      25,
		Password: "SecurePass123!",
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid combined example passed")
	}

	invalid := CombinedExample{
		Username: "ab",      // Too short (min: 3)
		Email:    "invalid", // Invalid email
		Age:      15,        // Too young (min: 18)
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught multiple validation errors (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateNestedStructures() {
	fmt.Println("6. NESTED STRUCTURES")
	fmt.Println("--------------------")

	valid := Person{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			ZipCode: "10001",
			Country: "US",
		},
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid nested structure example passed")
	}

	invalid := Person{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Address: Address{
			Street:  "", // Missing required field
			City:    "NY",
			ZipCode: "123",
			Country: "INVALID", // Not in allowed values
		},
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught nested validation errors (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateOptionalFields() {
	fmt.Println("7. OPTIONAL FIELDS (POINTERS)")
	fmt.Println("----------------------------")

	email := "test@example.com"
	number := 50

	valid := OptionalExample{
		RequiredField:  "required",
		OptionalString: &email,
		OptionalNumber: &number,
		OptionalEmail:  &email,
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid optional fields example passed")
	}

	// Valid with nil optional fields
	valid2 := OptionalExample{
		RequiredField: "required",
		// All optional fields are nil - this is valid
	}

	err = validator.Validate(valid2)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid with nil optional fields")
	}

	// Invalid - optional field provided but doesn't meet constraints
	shortString := "abc" // Too short (min: 5)
	invalid := OptionalExample{
		RequiredField:  "required",
		OptionalString: &shortString,
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught optional field validation error (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateSlices() {
	fmt.Println("8. SLICES AND ARRAYS")
	fmt.Println("-------------------")

	valid := SliceExample{
		RequiredSlice: []string{"item1", "item2"},
		Tags:          []string{"tag1", "tag2"},
		Numbers:       []int{1, 2, 3},
		OptionalSlice: []string{"optional"},
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid slice example passed")
	}

	invalid := SliceExample{
		RequiredSlice: []string{}, // Empty - violates required
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught slice validation error (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateMaps() {
	fmt.Println("9. MAPS")
	fmt.Println("-------")

	valid := MapExample{
		RequiredMap: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		ConfigMap: map[string]string{
			"port":    "8080",
			"timeout": "30",
		},
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid map example passed")
	}

	invalid := MapExample{
		RequiredMap: map[string]string{}, // Empty - violates required
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught map validation error (expected): %v\n", err)
	}
	fmt.Println()
}

func demonstrateRealWorldExample() {
	fmt.Println("10. REAL-WORLD EXAMPLE")
	fmt.Println("---------------------")

	valid := UserRegistration{
		Username:    "johndoe123",
		Email:       "john.doe@example.com",
		Password:    "SecurePass123!",
		ConfirmPass: "SecurePass123!",
		FirstName:   "John",
		LastName:    "Doe",
		Age:         25,
		Phone:       "1234567890",
		Website:     "https://johndoe.com",
		Status:      "active",
		Country:     "US",
		StartDate:   "2024-01-01",
		EndDate:     "2024-12-31",
		Address: Address{
			Street:  "123 Main Street",
			City:    "New York",
			ZipCode: "10001",
			Country: "US",
		},
		Tags: []string{"developer", "golang"},
	}

	err := validator.Validate(valid)
	if err != nil {
		fmt.Printf("❌ Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Valid user registration example passed")
	}

	invalid := UserRegistration{
		Username:    "ab",            // Too short
		Email:       "invalid-email", // Invalid format
		Password:    "short",         // Too short
		ConfirmPass: "short",
		FirstName:   "J",          // Too short
		Age:         15,           // Too young
		Status:      "invalid",    // Not in allowed values
		StartDate:   "2024-12-31", // After end date
		EndDate:     "2024-01-01",
	}

	err = validator.Validate(invalid)
	if err != nil {
		fmt.Printf("✓ Caught multiple validation errors (expected):\n%v\n", err)
	}
	fmt.Println()

	fmt.Println("=== All Examples Completed ===")
}

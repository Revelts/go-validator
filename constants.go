package validator

const (
	// Tag values
	tagRequired     = "required"
	tagEmail        = "email"
	tagAlphanumeric = "alphanumeric"
	tagAlphabet     = "alphabet"
	tagNumber       = "number"
	tagText         = "text"

	// Regex patterns (using different names to avoid conflicts with model.go)
	mailFormatRegex         = `\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`
	alphanumericFormatRegex = `^[A-Za-z0-9 ]*$`
	alphabetFormatRegex     = `^[A-Za-z ]*$`
	numericFormatRegex      = `^[0-9]*$`
)

package validator

// formatError is a simple error type to avoid linter warnings about non-constant format strings
type formatError struct {
	message string
}

func (e *formatError) Error() string {
	return e.message
}


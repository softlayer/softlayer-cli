package errors

import (
	"strings"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

// deprecated. Use InvalidUsageError instead
var ErrInvalidUsage = new(InvalidUsageError)

type InvalidUsageError struct {
	Message string
}

func NewInvalidUsageError(message string) *InvalidUsageError {
	return &InvalidUsageError{Message: message}
}

func (e *InvalidUsageError) Error() string {
	if e.Message == "" {
		return T("Incorrect Usage.")
	}

	return T("Incorrect Usage: ") + e.Message
}

// NewExclusiveFlagsError creates an InvalidUsageError about exclusive flags
func NewExclusiveFlagsError(flag1 string, flag2 string, moreFlags ...string) *InvalidUsageError {
	flags := []string{flag1, flag2}
	flags = append(flags, moreFlags...)
	return NewExclusiveFlagsErrorWithDetails(flags, "")
}

// NewExclusiveFlagsErrorWithDetails creates an InvalidUsageError about exclusive flags with details
func NewExclusiveFlagsErrorWithDetails(flags []string, details string) *InvalidUsageError {
	str := strings.Join(flags, "', '")
	str = "'" + str + "'"

	msg := T("{{.Flags}} are exclusive.", map[string]interface{}{
		"Flags": str,
	})

	if details != "" {
		msg = msg + " " + details
	}
	return &InvalidUsageError{
		Message: msg,
	}
}

// NewMissingInputError create an InvalidUsageError about a required input
func NewMissingInputError(inputName string) *InvalidUsageError {
	return &InvalidUsageError{
		Message: T("'{{.Input}}' is required", map[string]interface{}{
			"Input": inputName,
		}),
	}
}

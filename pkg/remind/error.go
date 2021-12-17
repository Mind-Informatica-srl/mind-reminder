package remind

import "fmt"

// RemindGenericError is raised on error
type RemindGenericError struct {
	Reason string `json:"reason"`
}

// NewRemindGenericError restituisce nuovo NewRemindGenericError passando una stringa in input
func NewRemindGenericError(reason string) RemindGenericError {
	return RemindGenericError{
		Reason: reason,
	}
}

// Error serve perchè RemindGenericError implementi Errors
func (a RemindGenericError) Error() string {
	return fmt.Sprintf("remind error: %v", a.Reason)
}

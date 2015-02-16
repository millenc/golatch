package golatch

import "fmt"

type LatchError struct {
	Code    string
	Message string
}

//Implementation of the error interface
func (e *LatchError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.Code, e.Message)
}

//Constructs a new error
func (e *LatchError) NewLatchError(code, message string) error {
	return &LatchError{Code: code,
		Message: message,
	}
}

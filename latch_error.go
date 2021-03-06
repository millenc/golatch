package golatch

import "fmt"

type LatchError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

//Implementation of the error interface
func (e *LatchError) Error() string {
	return fmt.Sprintf("Latch Error: [%d] %s", e.Code, e.Message)
}

//Constructs a new error
func NewLatchError(code int32, message string) error {
	return &LatchError{Code: code, Message: message}
}

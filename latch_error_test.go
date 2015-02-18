package golatch

import (
	"testing"
)

func TestNewLatchError(t *testing.T) {
	err := NewLatchError(205, "Account and application already paired")

	latch_error := err.(*LatchError)
	if latch_error.Code != 205 || latch_error.Message != "Account and application already paired" {
		t.Errorf("NewLatchError() failed")
	}
}

func TestLatchErrorError(t *testing.T) {
	err := &LatchError{Code: 205, Message: "Account and application already paired"}

	if err.Error() != "Latch Error: [205] Account and application already paired" {
		t.Errorf("NewLatchError() failed")
	}
}

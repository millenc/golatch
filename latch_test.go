package golatch

import (
	"fmt"
	"testing"
)

func TestNewLatch(t *testing.T) {
	latch := NewLatch("MyAppID", "MySecretKey")

	if latch.AppID != "MyAppID" || latch.SecretKey != "MySecretKey" {
		t.Errorf("NewLatch() failed: expected (%q,%q), got (%q,%q)", "MyAppID", "MySecretKey", latch.AppID, latch.SecretKey)
	}
}

func TestGetLatchQueryString(t *testing.T) {
	expected := "/api/0.9/pair/my_token"

	if got := GetLatchQueryString(fmt.Sprint(API_PAIR_ACTION, "/", "my_token")); got != expected {
		t.Errorf("GetLatchQueryString() failed: expected %q, got %q", expected, got)
	}
}

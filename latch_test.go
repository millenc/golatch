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

func TestGetLatchURL(t *testing.T) {
	expected_url := "https://latch.elevenpaths.com/api/1.0/pair/my_token"

	if got_url := GetLatchURL(fmt.Sprint(API_PAIR_ACTION, "/", "my_token")); got_url.String() != "https://latch.elevenpaths.com/api/1.0/pair/my_token" {
		t.Errorf("GetLatchURL() failed: expected %q, got %q", expected_url, got_url)
	}
}

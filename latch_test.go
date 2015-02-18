package golatch

import (
	"fmt"
	"testing"
)

func TestSetAppID(t *testing.T) {
	l := &Latch{}
	appID := "MyAppID"

	l.SetAppID(appID)
	if l.AppID() != appID {
		t.Errorf("SetAppID()/AppID() failed: expected %q, got %q", appID, l.AppID())
	}
}

func TestSetSecretKey(t *testing.T) {
	l := &Latch{}
	secretKey := "MySecretKey"

	l.SetSecretKey(secretKey)
	if l.SecretKey() != secretKey {
		t.Errorf("SetSecretKey()/SecretKey() failed: expected %q, got %q", secretKey, l.SecretKey())
	}
}

func TestGetLatchQueryString(t *testing.T) {
	expected := "/api/0.9/pair/my_token"

	if got := GetLatchQueryString(fmt.Sprint(API_PAIR_ACTION, "/", "my_token")); got != expected {
		t.Errorf("GetLatchQueryString() failed: expected %q, got %q", expected, got)
	}
}

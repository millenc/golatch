package golatch

import (
	"reflect"
	"testing"
	"time"
)

var example_date = time.Date(2015, time.February, 15, 14, 53, 0, 0, time.UTC) //2015-02-15 14:53:00
var example_request = &LatchRequest{AppID: "MyAppID",
	SecretKey:   "MySecretKey",
	HttpMethod:  "POST",
	QueryString: "/api/0.9/pair/my_token",
	XHeaders: map[string]string{
		"X-11paths-B": "Test value",
		"X-11paths-A": "Line\nBreaks",
	},
	Params: map[string][]string{
		"B[]": {"B", "A"},
		"A":   {"A"},
	},
	Date: example_date,
}
var example_expected_signature = "POST\n" +
	"2015-02-15 14:53:00\n" +
	"x-11paths-a:Line Breaks x-11paths-b:Test value\n" +
	"/api/0.9/pair/my_token\n" +
	"A=A&B[]=A&B[]=B"
var example_expected_header = "11PATHS MyAppID dKuWhQ2YCcx3c92bus5Bp6wo7kk="

func TestNewLatchRequest(t *testing.T) {
	got_request := NewLatchRequest(example_request.AppID, example_request.SecretKey, example_request.HttpMethod, example_request.QueryString, example_request.XHeaders, example_request.Params, example_date)

	if example_request.AppID != got_request.AppID ||
		example_request.SecretKey != got_request.SecretKey ||
		example_request.HttpMethod != got_request.HttpMethod ||
		example_request.QueryString != got_request.QueryString ||
		!reflect.DeepEqual(example_request.XHeaders, got_request.XHeaders) ||
		!reflect.DeepEqual(example_request.Params, got_request.Params) ||
		example_request.Date != got_request.Date {
		t.Errorf("NewLatchRequest() failed: expected %q, got %q", example_request, got_request)
	}
}

func TestGetAuthenticationHeaders(t *testing.T) {
	headers := example_request.GetAuthenticationHeaders()

	if headers == nil {
		t.Errorf("GetAuthenticationHeaders() failed: returned value should not be nil")
	}
	if _, ok := headers[API_AUTHORIZATION_HEADER_NAME]; !ok {
		t.Errorf("GetAuthenticationHeaders() failed: missing authorization header, got %q", headers)
	}
	if _, ok := headers[API_DATE_HEADER_NAME]; !ok {
		t.Errorf("GetAuthenticationHeaders() failed: missing authorization header date, got %q", headers)
	}
}

func TestGetRequestSignature(t *testing.T) {
	got_signature := example_request.GetRequestSignature()

	if example_expected_signature != got_signature {
		t.Errorf("GetRequestSignature() failed: expected %q (%x), got %q(%x)", example_expected_signature, example_expected_signature, got_signature, got_signature)
	}
}

func TestGetAuthorizationHeader(t *testing.T) {
	got_header := example_request.GetAuthorizationHeader()

	if example_expected_header != got_header {
		t.Errorf("GetAuthorizationHeader() failed: expected %q, got %q", example_expected_header, got_header)
	}
}

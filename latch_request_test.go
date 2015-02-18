package golatch

import (
	"reflect"
	"testing"
	"time"
)

//Example request with dummy data used on the following tests
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
var example_expected_formatted_date = "2015-02-15 14:53:00"
var example_expected_serialized_headers = "x-11paths-a:Line Breaks x-11paths-b:Test value"
var example_expected_serialized_params = "A=A&B[]=A&B[]=B"

var example_expected_signature = "POST\n" +
	example_expected_formatted_date + "\n" +
	example_expected_serialized_headers + "\n" +
	"/api/0.9/pair/my_token\n" +
	example_expected_serialized_params
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
	if _, ok := headers["Authorization"]; !ok {
		t.Errorf("GetAuthenticationHeaders() failed: missing authorization header, got %q", headers)
	}
	if _, ok := headers["X-11Paths-Date"]; !ok {
		t.Errorf("GetAuthenticationHeaders() failed: missing authorization header date, got %q", headers)
	}
}

func TestGetRequestSignature(t *testing.T) {
	got_signature := example_request.GetRequestSignature()

	if example_expected_signature != got_signature {
		t.Errorf("GetRequestSignature() failed: expected %q (%x), got %q(%x)", example_expected_signature, example_expected_signature, got_signature, got_signature)
	}
}

func TestGetSerializedHeaders(t *testing.T) {
	got_headers := example_request.GetSerializedHeaders()

	if example_expected_serialized_headers != got_headers {
		t.Errorf("GetSerializedHeaders() failed: expected %q, got %q", example_expected_serialized_headers, got_headers)
	}
}

func TestGetSerializedParams(t *testing.T) {
	got_params := example_request.GetSerializedParams()

	if example_expected_serialized_params != got_params {
		t.Errorf("GetSerializedParams() failed: expected %q, got %q", example_expected_serialized_params, got_params)
	}
}

func TestGetAuthorizationHeader(t *testing.T) {
	got_header := example_request.GetAuthorizationHeader()

	if example_expected_header != got_header {
		t.Errorf("GetAuthorizationHeader() failed: expected %q, got %q", example_expected_header, got_header)
	}
}

func TestGetFormattedDate(t *testing.T) {
	got_formatted_date := example_request.GetFormattedDate()

	if example_expected_formatted_date != got_formatted_date {
		t.Errorf("GetFormattedDate() failed: expected %q, got %q", example_expected_formatted_date, got_formatted_date)
	}
}

func TestGetHttpRequest(t *testing.T) {
	got_request := example_request.GetHttpRequest()

	if got_request.Method != example_request.HttpMethod {
		t.Errorf("GetHttpRequest() failed: expected HTTP Method %q, got %q", example_request.HttpMethod, got_request.Method)
	}
	if got_request.URL.String() != "https://latch.elevenpaths.com/api/0.9/pair/my_token" {
		t.Errorf("GetHttpRequest() failed: expected URL %q, got %q", "https://latch.elevenpaths.com/api/0.9/pair/my_token", got_request.URL.String())
	}
	if got_request.Header.Get(API_AUTHORIZATION_HEADER_NAME) != example_request.GetAuthorizationHeader() {
		t.Errorf("GetHttpRequest() failed: expected Authorization header %q, got %q", got_request.Header.Get(API_AUTHORIZATION_HEADER_NAME), example_request.GetAuthorizationHeader())
	}
	if got_request.Header.Get(API_DATE_HEADER_NAME) != example_request.GetFormattedDate() {
		t.Errorf("GetHttpRequest() failed: expected Date header %q, got %q", got_request.Header.Get(API_DATE_HEADER_NAME), example_request.GetFormattedDate())
	}
}

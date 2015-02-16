package golatch

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type LatchRequest struct {
	AppID       string
	SecretKey   string
	HttpMethod  string
	QueryString string
	XHeaders    map[string]string
	Params      map[string][]string
	Date        time.Time
}

//Returns a new LatchRequest initialized with the parameters provided
func NewLatchRequest(appID string, secretKey string, httpMethod string, queryString string, xHeaders map[string]string, params map[string][]string, date time.Time) *LatchRequest {
	return &LatchRequest{
		AppID:       appID,
		SecretKey:   secretKey,
		HttpMethod:  httpMethod,
		QueryString: queryString,
		XHeaders:    xHeaders,
		Params:      params,
		Date:        date,
	}
}

//Gets the authentication headers (Authorization and Date)
func (l *LatchRequest) GetAuthenticationHeaders() (headers map[string]string) {
	headers = make(map[string]string)
	headers[API_AUTHORIZATION_HEADER_NAME] = l.GetAuthorizationHeader()
	headers[API_DATE_HEADER_NAME] = l.GetFormattedDate()

	return headers
}

//Gets the Authorization header
func (l *LatchRequest) GetAuthorizationHeader() string {
	return fmt.Sprint(API_AUTHENTICATION_METHOD,
		API_AUTHORIZATION_HEADER_FIELD_SEPARATOR,
		l.AppID,
		API_AUTHORIZATION_HEADER_FIELD_SEPARATOR,
		l.GetSignedRequestSignature())
}

//Gets the signed request signature using HMAC-SHA1 (base64-encoded)
func (l *LatchRequest) GetSignedRequestSignature() string {
	key := []byte(l.SecretKey)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(l.GetRequestSignature()))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

//Gets the request signature
func (l *LatchRequest) GetRequestSignature() string {
	var buffer bytes.Buffer

	buffer.WriteString(l.HttpMethod)
	buffer.WriteString("\n")
	buffer.WriteString(l.GetFormattedDate())
	buffer.WriteString("\n")
	buffer.WriteString(l.GetSerializedHeaders())
	buffer.WriteString("\n")
	buffer.WriteString(l.QueryString)
	if l.HttpMethod == HTTP_METHOD_POST || l.HttpMethod == HTTP_METHOD_PUT {
		buffer.WriteString("\n")
		buffer.WriteString(l.GetSerializedParams())
	}

	return buffer.String()
}

//Gets the serialized request headers (xHeaders) to use in the request signature
func (l *LatchRequest) GetSerializedHeaders() string {
	if l.XHeaders == nil {
		return ""
	}

	//Get lowercase header names and remove line breaks
	m := make(map[string]string)
	var keys []string
	for k, v := range l.XHeaders {
		var key string = strings.ToLower(k)
		keys = append(keys, key)
		m[key] = strings.Replace(v, "\n", " ", -1)
	}

	//Sort header names in ascending alphabetical order
	sort.Strings(keys)

	//Serialize
	var buffer bytes.Buffer
	for _, key := range keys {
		buffer.WriteString(fmt.Sprint(key, ":", m[key], " "))
	}

	return strings.Trim(buffer.String(), " ")
}

//Gets the serialized request headers (xHeaders) to use in the request signature
func (l *LatchRequest) GetSerializedParams() string {
	if l.Params == nil {
		return ""
	}

	//Get parameter names and sort them alphabetically
	m := make(map[string][]string)
	var keys []string
	for k, v := range l.Params {
		keys = append(keys, k)
		sort.Strings(v) //Sort parameter values
		m[k] = v
	}
	sort.Strings(keys)

	//Serialize
	var buffer bytes.Buffer
	for _, key := range keys {
		values := m[key]

		for _, value := range values {
			buffer.WriteString(fmt.Sprint(key, "=", value, "&"))
		}
	}

	return strings.Trim(buffer.String(), " &")
}

//Gets the current UTC Date/Time as a string formatted using the layout specified in const(API_UTC_STRING_FORMAT)
func (l *LatchRequest) GetFormattedDate() string {
	return l.Date.UTC().Format(API_UTC_STRING_FORMAT)
}

//Gets the HTTP request for this Latch Request
func (l *LatchRequest) GetHttpRequest() *http.Request {
	//TODO: Get the URL from the LatchRequest?
	request, _ := http.NewRequest(l.HttpMethod, fmt.Sprint(API_URL, l.QueryString), nil)

	//Set Headers
	headers := l.GetAuthenticationHeaders()
	for header, value := range headers {
		request.Header.Set(header, value)
	}

	return request
}

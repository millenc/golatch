package golatch

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

type LatchRequest struct {
	AppID       string
	HttpMethod  string
	QueryString string
	XHeaders    map[string]string
	Params      map[string]string
	Date        time.Time
}

//Returns a new LatchRequest initialized with the parameters provided
func NewLatchRequest(appID string, httpMethod string, queryString string, xHeaders map[string]string, params map[string]string, date time.Time) *LatchRequest {
	return &LatchRequest{
		AppID:       appID,
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
		l.GetRequestSignature())
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
	return ""
}

//Gets the current UTC Date/Time as a string formatted using the layout specified in const(API_UTC_STRING_FORMAT)
func (l *LatchRequest) GetFormattedDate() string {
	return l.Date.UTC().Format(API_UTC_STRING_FORMAT)
}

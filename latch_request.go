package golatch

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type LatchRequest struct {
	AppID      string
	SecretKey  string
	HttpMethod string
	URL        *url.URL
	XHeaders   map[string]string
	Params     url.Values
	Date       time.Time
}

//Returns a new LatchRequest initialized with the parameters provided
func NewLatchRequest(appID string, secretKey string, httpMethod string, url *url.URL, xHeaders map[string]string, params url.Values, date time.Time) *LatchRequest {
	return &LatchRequest{
		AppID:      appID,
		SecretKey:  secretKey,
		HttpMethod: httpMethod,
		URL:        url,
		XHeaders:   xHeaders,
		Params:     params,
		Date:       date,
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
	buffer.WriteString(l.URL.RequestURI())
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

	//Order parameter values alphabetically. url.Encode() will take care of ordering the keys
	ordered_params := url.Values{}
	for param, values := range l.Params {
		sort.Strings(values) //Sort parameter values
		for _, value := range values {
			ordered_params.Add(param, value)
		}
	}

	return strings.Trim(ordered_params.Encode(), " &")
}

//Gets the current UTC Date/Time as a string formatted using the layout specified in const(API_UTC_STRING_FORMAT)
func (l *LatchRequest) GetFormattedDate() string {
	return l.Date.UTC().Format(API_UTC_STRING_FORMAT)
}

//Gets the HTTP request for this Latch Request
func (l *LatchRequest) GetHttpRequest() *http.Request {
	var body io.Reader = nil

	//Include parameters for POST and PUT methods
	if l.HttpMethod == HTTP_METHOD_PUT || l.HttpMethod == HTTP_METHOD_POST {
		body = strings.NewReader(l.Params.Encode())
	}

	request, _ := http.NewRequest(l.HttpMethod, l.URL.String(), body)

	//Set Headers
	headers := l.GetAuthenticationHeaders()
	for header, value := range headers {
		request.Header.Set(header, value)
	}
	if l.HttpMethod == HTTP_METHOD_PUT || l.HttpMethod == HTTP_METHOD_POST {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	request.Header.Set("User-Agent", HTTP_USER_AGENT)

	return request
}

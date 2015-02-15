package golatch

import (
	"fmt"
	t "time"
)

const (
	//Latch related constants
	API_VERSION                              = "0.9"
	API_URL                                  = "https://latch.elevenpaths.com/api"
	API_CHECK_STATUS_ACTION                  = "status"
	API_PAIR_ACTION                          = "pair"
	API_PAIR_WITH_ID_ACTION                  = "pairWithId"
	API_UNPAIR_ACTION                        = "unpair"
	API_LOCK_ACTION                          = "lock"
	API_UNLOCK_ACTION                        = "unlock"
	API_HISTORY_ACTION                       = "history"
	API_OPERATION_ACTION                     = "operation"
	API_AUTHENTICATION_METHOD                = "11PATHS"
	API_AUTHORIZATION_HEADER_NAME            = "Authorization"
	API_DATE_HEADER_NAME                     = "X-11Paths-Date"
	API_AUTHORIZATION_HEADER_FIELD_SEPARATOR = " "
	API_X_11PATHS_HEADER_PREFIX              = "X-11Paths-"
	API_X_11PATHS_HEADER_SEPARATOR           = ":"
	API_UTC_STRING_FORMAT                    = "2006-01-02 15:04:05" //format layout as defined here: http://golang.org/pkg/time/#pkg-constants

	//HTTP methods
	HTTP_METHOD_POST   = "POST"
	HTTP_METHOD_GET    = "GET"
	HTTP_METHOD_PUT    = "PUT"
	HTTP_METHOD_DELETE = "DELETE"
)

type Latch struct {
	appID     string
	secretKey string
}

//Gets the Application ID
func (l *Latch) AppID() string {
	return l.appID
}

//Sets the Application ID
func (l *Latch) SetAppID(appID string) {
	l.appID = appID
}

//Gets the Secret Key
func (l *Latch) SecretKey() string {
	return l.secretKey
}

//Sets the Secret Key
func (l *Latch) SetSecretKey(secretKey string) {
	l.secretKey = secretKey
}

//Get the authorization header
func (l *Latch) GetAuthorizationHeader() string {
	return fmt.Sprint(API_AUTHENTICATION_METHOD,
		API_AUTHORIZATION_HEADER_FIELD_SEPARATOR,
		l.AppID(),
		API_AUTHORIZATION_HEADER_FIELD_SEPARATOR,
		l.GetRequestSignature())
}

//Gets the request signature
func (l *Latch) GetRequestSignature() string {
	return ""
}

//Gets the current UTC Date/Time as a string formatted using the layout specified in const(API_UTC_STRING_FORMAT)
func (l *Latch) getCurrentDateTime() string {
	return t.Now().UTC().Format(API_UTC_STRING_FORMAT)
}

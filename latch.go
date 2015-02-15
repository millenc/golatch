package golatch

import (
	"fmt"
	t "time"
)

const (
	//Latch related constants
	API_URL                                  = "https://latch.elevenpaths.com"
	API_PATH                                 = "/api"
	API_VERSION                              = "0.9"
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

//Pairs an account with the pairing token
func (l *Latch) Pair(token string) *LatchRequest {
	request := NewLatchRequest(l.AppID(), HTTP_METHOD_GET, GetLatchQueryString(fmt.Sprint(API_PAIR_ACTION, "/", token)), nil, nil, t.Now())
	return request
}

//Gets the complete url for a request
func GetLatchQueryString(queryString string) string {
	return fmt.Sprint(API_PATH, "/", API_VERSION, "/", queryString)
}

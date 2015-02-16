package golatch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
func (l *Latch) Pair(token string) (*LatchPairResponse, error) {
	response, err := l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_GET, GetLatchQueryString(fmt.Sprint(API_PAIR_ACTION, "/", token)), nil, nil, t.Now()), &LatchPairResponse{})
	return (*response).(*LatchPairResponse), err
}

func (l *Latch) DoRequest(request *LatchRequest, responseType LatchResponse) (response *LatchResponse, err error) {
	client := &http.Client{}
	req := request.GetHttpRequest()

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	//Get the response's body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	//Handle HTTP errors
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("HTTP error [%d] body: %s", resp.StatusCode, body))
		return
	}

	//TODO: Check if the response is an error before decoding it

	//Decode response into a typed response (if one has been specified)
	if responseType != nil {
		err = responseType.Unmarshal(string(body))
		response = &responseType
	}

	return response, err
}

//Gets the complete url for a request
func GetLatchQueryString(queryString string) string {
	return fmt.Sprint(API_PATH, "/", API_VERSION, "/", queryString)
}

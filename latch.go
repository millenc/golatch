package golatch

import (
	"encoding/json"
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
	API_NOOTP_SUFFIX                         = "nootp"
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

//Pairs an account with the provided pairing token
func (l *Latch) Pair(token string) (response *LatchPairResponse, err error) {
	var resp *LatchResponse
	if resp, err = l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_GET, GetLatchQueryString(fmt.Sprint(API_PAIR_ACTION, "/", token)), nil, nil, t.Now()), &LatchPairResponse{}); err == nil {
		response = (*resp).(*LatchPairResponse)
	}
	return response, err
}

//Unpairs an account, given it's account ID
func (l *Latch) Unpair(accountId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_GET, GetLatchQueryString(fmt.Sprint(API_UNPAIR_ACTION, "/", accountId)), nil, nil, t.Now()), nil)
	return err
}

//Locks an account, given it's account ID
func (l *Latch) Lock(accountId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_GET, GetLatchQueryString(fmt.Sprint(API_LOCK_ACTION, "/", accountId)), nil, nil, t.Now()), nil)
	return err
}

//Unlocks an account, given it's account ID
func (l *Latch) Unlock(accountId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_GET, GetLatchQueryString(fmt.Sprint(API_UNLOCK_ACTION, "/", accountId)), nil, nil, t.Now()), nil)
	return err
}

//Locks an operation, given it's account ID and oeration ID
func (l *Latch) LockOperation(accountId string, operationId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_GET, GetLatchQueryString(fmt.Sprint(API_LOCK_ACTION, "/", accountId, "/op/", operationId)), nil, nil, t.Now()), nil)
	return err
}

//Unlocks an operation, given it's account ID and oeration ID
func (l *Latch) UnlockOperation(accountId string, operationId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_GET, GetLatchQueryString(fmt.Sprint(API_UNLOCK_ACTION, "/", accountId, "/op/", operationId)), nil, nil, t.Now()), nil)
	return err
}

//Adds a new operation
func (l *Latch) AddOperation(parentId string, name string, twoFactor string, lockOnRequest string) (response *LatchOperationResponse, err error) {
	var resp *LatchResponse

	params := map[string][]string{
		"parentId":        {parentId},
		"name":            {name},
		"two_factor":      {twoFactor},
		"lock_on_request": {lockOnRequest},
	}

	if resp, err = l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_PUT, GetLatchQueryString(API_OPERATION_ACTION), nil, params, t.Now()), &LatchOperationResponse{}); err == nil {
		response = (*resp).(*LatchOperationResponse)
	}
	return response, err
}

//Gets the status of an account, given it's account ID
//If nootp is true, the one time password won't be included in the response
func (l *Latch) Status(accountId string, nootp bool) (response *LatchStatusResponse, err error) {
	query := fmt.Sprint(API_CHECK_STATUS_ACTION, "/", accountId)
	if nootp {
		query = fmt.Sprint(query, "/nootp")
	}

	return l.StatusRequest(query)
}

//Gets the status of an operation, given it's account ID and operation ID
//If nootp is true, the one time password won't be included in the response
func (l *Latch) OperationStatus(accountId string, operationId string, nootp bool) (response *LatchStatusResponse, err error) {
	query := fmt.Sprint(API_CHECK_STATUS_ACTION, "/", accountId, "/op/", operationId)
	if nootp {
		query = fmt.Sprint(query, "/nootp")
	}

	return l.StatusRequest(query)
}

//Performs a status request (application or operation) against the query URL provided
//Returns a LatchStatusResponse struct on success
func (l *Latch) StatusRequest(query string) (response *LatchStatusResponse, err error) {
	var resp *LatchResponse
	if resp, err = l.DoRequest(NewLatchRequest(l.AppID(), l.SecretKey(), HTTP_METHOD_GET, GetLatchQueryString(query), nil, nil, t.Now()), &LatchStatusResponse{}); err == nil {
		response = (*resp).(*LatchStatusResponse)
	}
	return response, err
}

func (l *Latch) DoRequest(request *LatchRequest, responseType LatchResponse) (response *LatchResponse, err error) {
	client := &http.Client{}
	req := request.GetHttpRequest()
	var resp *http.Response
	var body []byte

	//Perform the request
	if resp, err = client.Do(req); err != nil {
		return
	}

	//Get the response's body
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	//Handle HTTP errors
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("HTTP error [%d] body: %s", resp.StatusCode, body))
		return
	}

	//Check if the response is an error before decoding it
	latch_error_response := &LatchErrorResponse{}
	if err = json.Unmarshal(body, latch_error_response); err != nil {
		return
	} else if (*latch_error_response).Err.Code != 0 {
		err = &latch_error_response.Err
		return
	}

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

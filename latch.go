package golatch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	t "time"
)

const (
	//Latch related constants
	API_URL                                  = "https://latch.elevenpaths.com"
	API_PATH                                 = "/api"
	API_VERSION                              = "1.0"
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

	//Possible values for the Two factor and Lock on request options
	//NOT_SET is used in the UpdateOperation() method to leave the existing value
	MANDATORY = "MANDATORY"
	OPT_IN    = "OPT_IN"
	DISABLED  = "DISABLED"
	NOT_SET   = ""

	//Possible status values for the latch
	LATCH_STATUS_ON  = "on"
	LATCH_STATUS_OFF = "off"

	//HTTP methods
	HTTP_METHOD_POST   = "POST"
	HTTP_METHOD_GET    = "GET"
	HTTP_METHOD_PUT    = "PUT"
	HTTP_METHOD_DELETE = "DELETE"
)

type Latch struct {
	AppID     string
	SecretKey string
}

//Constructs a new Latch struct
func NewLatch(appID string, secretKey string) *Latch {
	return &Latch{
		AppID:     appID,
		SecretKey: secretKey,
	}
}

//Pairs an account with the provided pairing token
func (l *Latch) Pair(token string) (response *LatchPairResponse, err error) {
	var resp *LatchResponse
	if resp, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(fmt.Sprint(API_PAIR_ACTION, "/", token)), nil, nil, t.Now()), &LatchPairResponse{}); err == nil {
		response = (*resp).(*LatchPairResponse)
	}
	return response, err
}

//Unpairs an account, given it's account ID
func (l *Latch) Unpair(accountId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(fmt.Sprint(API_UNPAIR_ACTION, "/", accountId)), nil, nil, t.Now()), nil)
	return err
}

//Locks an account, given it's account ID
func (l *Latch) Lock(accountId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(fmt.Sprint(API_LOCK_ACTION, "/", accountId)), nil, nil, t.Now()), nil)
	return err
}

//Unlocks an account, given it's account ID
func (l *Latch) Unlock(accountId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(fmt.Sprint(API_UNLOCK_ACTION, "/", accountId)), nil, nil, t.Now()), nil)
	return err
}

//Locks an operation, given it's account ID and oeration ID
func (l *Latch) LockOperation(accountId string, operationId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(fmt.Sprint(API_LOCK_ACTION, "/", accountId, "/op/", operationId)), nil, nil, t.Now()), nil)
	return err
}

//Unlocks an operation, given it's account ID and oeration ID
func (l *Latch) UnlockOperation(accountId string, operationId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(fmt.Sprint(API_UNLOCK_ACTION, "/", accountId, "/op/", operationId)), nil, nil, t.Now()), nil)
	return err
}

//Adds a new operation
func (l *Latch) AddOperation(parentId string, name string, twoFactor string, lockOnRequest string) (response *LatchAddOperationResponse, err error) {
	var resp *LatchResponse

	params := url.Values{}
	params.Set("parentId", parentId)
	params.Set("name", name)
	params.Set("two_factor", twoFactor)
	params.Set("lock_on_request", lockOnRequest)

	if resp, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_PUT, GetLatchURL(API_OPERATION_ACTION), nil, params, t.Now()), &LatchAddOperationResponse{}); err == nil {
		response = (*resp).(*LatchAddOperationResponse)
	}
	return response, err
}

//Updates an existing operation
func (l *Latch) UpdateOperation(operationId string, name string, twoFactor string, lockOnRequest string) (err error) {
	params := url.Values{}
	params.Set("name", name)
	if twoFactor != NOT_SET {
		params.Set("two_factor", twoFactor)
	}
	if lockOnRequest != NOT_SET {
		params.Set("lock_on_request", lockOnRequest)
	}

	_, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_POST, GetLatchURL(fmt.Sprint(API_OPERATION_ACTION, "/", operationId)), nil, params, t.Now()), nil)
	return err
}

//Deletes an existing operation
func (l *Latch) DeleteOperation(operationId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_DELETE, GetLatchURL(fmt.Sprint(API_OPERATION_ACTION, "/", operationId)), nil, nil, t.Now()), nil)
	return err
}

//Shows operations information
//If operationId is empty this function will retrieve all the operations of the app
func (l *Latch) ShowOperation(operationId string) (response *LatchShowOperationResponse, err error) {
	var resp *LatchResponse
	var operation string

	if operationId != "" {
		operation += "/" + operationId
	}

	if resp, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(fmt.Sprint(API_OPERATION_ACTION, operation)), nil, nil, t.Now()), &LatchShowOperationResponse{}); err == nil {
		response = (*resp).(*LatchShowOperationResponse)
	}
	return response, err
}

//Gets the status of an account, given it's account ID
//If nootp is true, the one time password won't be included in the response
//If silent is true Latch will not send push notifications to the client (requires SILVER, GOLD or PLATINUM subscription)
func (l *Latch) Status(accountId string, nootp bool, silent bool) (response *LatchStatusResponse, err error) {
	query := fmt.Sprint(API_CHECK_STATUS_ACTION, "/", accountId)
	if nootp {
		query = fmt.Sprint(query, "/nootp")
	}
	if silent {
		query = fmt.Sprint(query, "/silent")
	}

	return l.StatusRequest(query)
}

//Gets the status of an operation, given it's account ID and operation ID
//If nootp is true, the one time password won't be included in the response
//If silent is true Latch will not send push notifications to the client (requires SILVER, GOLD or PLATINUM subscription)
func (l *Latch) OperationStatus(accountId string, operationId string, nootp bool, silent bool) (response *LatchStatusResponse, err error) {
	query := fmt.Sprint(API_CHECK_STATUS_ACTION, "/", accountId, "/op/", operationId)
	if nootp {
		query = fmt.Sprint(query, "/nootp")
	}
	if silent {
		query = fmt.Sprint(query, "/silent")
	}

	return l.StatusRequest(query)
}

//Performs a status request (application or operation) against the query URL provided
//Returns a LatchStatusResponse struct on success
func (l *Latch) StatusRequest(query string) (response *LatchStatusResponse, err error) {
	var resp *LatchResponse
	if resp, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(query), nil, nil, t.Now()), &LatchStatusResponse{}); err == nil {
		response = (*resp).(*LatchStatusResponse)
	}
	return response, err
}

//Gets the account's history between the from and to dates
func (l *Latch) History(accountId string, from t.Time, to t.Time) (response *LatchHistoryResponse, err error) {
	var resp *LatchResponse

	if resp, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(fmt.Sprintf("%s/%s/%d/%d", API_HISTORY_ACTION, accountId, from.UnixNano()/1000000, to.UnixNano()/1000000)), nil, nil, t.Now()), &LatchHistoryResponse{AppID: l.AppID}); err == nil {
		response = (*resp).(*LatchHistoryResponse)
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
func GetLatchURL(queryString string) *url.URL {
	latch_url, err := (&url.URL{}).Parse(fmt.Sprint(API_URL, API_PATH, "/", API_VERSION, "/", queryString))
	if err != nil {
		latch_url = &url.URL{}
	}

	return latch_url
}

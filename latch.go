package golatch

import (
	"fmt"
	"net/url"
	t "time"
)

type Latch struct {
	AppID     string
	SecretKey string
	LatchAPI
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
		query = fmt.Sprint(query, "/", API_NOOTP_SUFFIX)
	}
	if silent {
		query = fmt.Sprint(query, "/", API_SILENT_SUFFIX)
	}

	return l.StatusRequest(query)
}

//Gets the status of an operation, given it's account ID and operation ID
//If nootp is true, the one time password won't be included in the response
//If silent is true Latch will not send push notifications to the client (requires SILVER, GOLD or PLATINUM subscription)
func (l *Latch) OperationStatus(accountId string, operationId string, nootp bool, silent bool) (response *LatchStatusResponse, err error) {
	query := fmt.Sprint(API_CHECK_STATUS_ACTION, "/", accountId, "/op/", operationId)
	if nootp {
		query = fmt.Sprint(query, "/", API_NOOTP_SUFFIX)
	}
	if silent {
		query = fmt.Sprint(query, "/", API_SILENT_SUFFIX)
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

	query := fmt.Sprintf("%s/%s", API_HISTORY_ACTION, accountId)
	if !from.IsZero() || !to.IsZero() {
		if !from.IsZero() {
			query = fmt.Sprint(query, fmt.Sprintf("/%d", from.UnixNano()/1000000))
		} else {
			query = fmt.Sprint(query, fmt.Sprintf("/%d", 0))
		}
	}
	if !to.IsZero() {
		query = fmt.Sprint(query, fmt.Sprintf("/%d", to.UnixNano()/1000000))
	}

	if resp, err = l.DoRequest(NewLatchRequest(l.AppID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(query), nil, nil, t.Now()), &LatchHistoryResponse{AppID: l.AppID}); err == nil {
		response = (*resp).(*LatchHistoryResponse)
	}
	return response, err
}

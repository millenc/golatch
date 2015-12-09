package golatch

import (
	"fmt"
	"net/url"
	t "time"
)

//Struct to use the Latch User API
type LatchUser struct {
	UserID    string
	SecretKey string
	LatchAPI
}

//Constructs a new LatchUser struct
func NewLatchUser(userID string, secretKey string) *LatchUser {
	return &LatchUser{
		UserID:    userID,
		SecretKey: secretKey,
	}
}

//Gets the user's subscription information
func (l *LatchUser) Subscription() (response *LatchSubscriptionResponse, err error) {
	var resp *LatchResponse
	if resp, err = l.DoRequest(NewLatchRequest(l.UserID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(API_SUBSCRIPTION_ACTION), nil, nil, t.Now()), &LatchSubscriptionResponse{}); err == nil {
		response = (*resp).(*LatchSubscriptionResponse)
	}

	return response, err
}

//Shows existing applications
func (l *LatchUser) ShowApplications() (response *LatchShowApplicationsResponse, err error) {
	var resp *LatchResponse
	if resp, err = l.DoRequest(NewLatchRequest(l.UserID, l.SecretKey, HTTP_METHOD_GET, GetLatchURL(API_APPLICATION_ACTION), nil, nil, t.Now()), &LatchShowApplicationsResponse{}); err == nil {
		response = (*resp).(*LatchShowApplicationsResponse)
	}
	return response, err
}

//Adds a new application
func (l *LatchUser) AddApplication(applicationInfo *LatchApplicationInfo) (response *LatchAddApplicationResponse, err error) {
	var resp *LatchResponse

	params := prepareApplicationParams(applicationInfo)

	if resp, err = l.DoRequest(NewLatchRequest(l.UserID, l.SecretKey, HTTP_METHOD_PUT, GetLatchURL(API_APPLICATION_ACTION), nil, *params, t.Now()), &LatchAddApplicationResponse{}); err == nil {
		response = (*resp).(*LatchAddApplicationResponse)
	}
	return response, err
}

//Updates application information
func (l *LatchUser) UpdateApplication(appID string, applicationInfo *LatchApplicationInfo) (err error) {
	params := prepareApplicationParams(applicationInfo)

	_, err = l.DoRequest(NewLatchRequest(l.UserID, l.SecretKey, HTTP_METHOD_POST, GetLatchURL(fmt.Sprint(API_APPLICATION_ACTION, "/", appID)), nil, *params, t.Now()), nil)

	return err
}

//Deletes an existing application
func (l *LatchUser) DeleteApplication(applicationId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.UserID, l.SecretKey, HTTP_METHOD_DELETE, GetLatchURL(fmt.Sprint(API_APPLICATION_ACTION, "/", applicationId)), nil, nil, t.Now()), nil)
	return err
}

//Initializes params for adding/updating application information
func prepareApplicationParams(applicationInfo *LatchApplicationInfo) (params *url.Values) {
	params = &url.Values{}
	params.Set("name", applicationInfo.Name)
	params.Set("contactEmail", applicationInfo.ContactEmail)
	params.Set("contactPhone", applicationInfo.ContactPhone)
	params.Set("two_factor", applicationInfo.TwoFactor)
	params.Set("lock_on_request", applicationInfo.LockOnRequest)

	return params
}

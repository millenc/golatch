package golatch

import (
	"fmt"
	"net/url"
	t "time"
)

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

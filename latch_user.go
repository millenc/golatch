package golatch

import (
	"fmt"
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

//Deletes an existing application
func (l *LatchUser) DeleteApplication(applicationId string) (err error) {
	_, err = l.DoRequest(NewLatchRequest(l.UserID, l.SecretKey, HTTP_METHOD_DELETE, GetLatchURL(fmt.Sprint(API_APPLICATION_ACTION, "/", applicationId)), nil, nil, t.Now()), nil)
	return err
}

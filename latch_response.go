package golatch

import (
	"encoding/json"
	"strings"
)

type LatchResponse interface {
	Unmarshal(Json string) error
}

type LatchErrorResponse struct {
	Err LatchError `json:"error"`
}

type LatchPairResponse struct {
	Data struct {
		AccountId string `json:"accountId"`
	} `json:"data"`
}

type LatchStatusResponse struct {
	Data struct {
		Operations map[string]LatchOperationStatus `json:"operations"`
	} `json:"data"`
}

type LatchOperationStatus struct {
	Status     string                          `json:"status"`
	TwoFactor  LatchTwoFactor                  `json:"two_factor"`
	Operations map[string]LatchOperationStatus `json:"operations"`
}

type LatchTwoFactor struct {
	Token     string `json:"token"`
	Generated int64  `json:"generated"`
}

type LatchAddOperationResponse struct {
	Data struct {
		OperationId string `json:"operationId"`
	} `json:"data"`
}

type LatchShowOperationResponse struct {
	Data struct {
		Operations map[string]LatchOperation `json:"operations"`
	} `json:"data"`
}

type LatchOperation struct {
	Name          string                    `json:"name"`
	Status        string                    `json:"status"`
	TwoFactor     string                    `json:"two_factor"`
	LockOnRequest string                    `json:"lock_on_request"`
	Operations    map[string]LatchOperation `json:"operations"`
}

type LatchHistoryResponse struct {
	AppID string
	Data  struct {
		Application   LatchApplication     `json:"application"`
		LastSeen      int64                `json:"lastSeen"`
		ClientVersion []LatchClientVersion `json:"clientVersion"`
		HistoryCount  int                  `json:"count"`
		History       []LatchHistoryEntry  `json:"history"`
	} `json:"data"`
}

type LatchApplication struct {
	Status   string `json:"status"`
	PairedOn int64  `json:"pairedOn"`
	LatchApplicationInfo
}

type LatchClientVersion struct {
	Platform string `json:platform`
	App      string `json:app`
}

type LatchHistoryEntry struct {
	Time      int64  `json:"t"`
	Action    string `json:"action"`
	What      string `json:"what"`
	Value     string `json:"value"`
	Was       string `json:"was"`
	Name      string `json:"name"`
	UserAgent string `json:"userAgent"`
	IP        string `json:"ip"`
}

type LatchShowApplicationsResponse struct {
	Data struct {
		Applications map[string]LatchApplicationInfo `json:"operations"`
	} `json:"data"`
}

type LatchApplicationInfo struct {
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	Secret        string                    `json:"secret"`
	ImageURL      string                    `json:"imageURL"`
	ContactPhone  string                    `json:"contactPhone"`
	ContactEmail  string                    `json:"contactEmail"`
	TwoFactor     string                    `json:"two_factor"`
	LockOnRequest string                    `json:"lock_on_request"`
	Operations    map[string]LatchOperation `json:"operations"`
}

type LatchAddApplicationResponse struct {
	Data struct {
		AppID  string `json:"applicationId"`
		Secret string `json:"secret"`
	} `json:"data"`
}

type LatchSubscriptionResponse struct {
	Data struct {
		Subscription struct {
			ID           string                            `json:"id"`
			Applications LatchSubscriptionUsage            `json:"applications"`
			Operations   map[string]LatchSubscriptionUsage `json:"operations"`
			Users        LatchSubscriptionUsage            `json:"users"`
		} `json:"subscription"`
	} `json:"data"`
}

type LatchSubscriptionUsage struct {
	InUse int `json:"inUse"`
	Limit int `json:"limit"`
}

func (l *LatchErrorResponse) Unmarshal(Json string) error {
	return json.Unmarshal([]byte(Json), l)
}

func (l *LatchPairResponse) Unmarshal(Json string) error {
	return json.Unmarshal([]byte(Json), l)
}

func (l *LatchStatusResponse) Unmarshal(Json string) (err error) {
	return json.Unmarshal([]byte(Json), l)
}

func (l *LatchAddOperationResponse) Unmarshal(Json string) (err error) {
	return json.Unmarshal([]byte(Json), l)
}

func (l *LatchShowOperationResponse) Unmarshal(Json string) (err error) {
	return json.Unmarshal([]byte(Json), l)
}

func (l *LatchHistoryResponse) Unmarshal(Json string) (err error) {
	return json.Unmarshal([]byte(strings.Replace(Json, l.AppID, "application", 1)), l)
}

func (l *LatchPairResponse) AccountId() string {
	return l.Data.AccountId
}

func (l *LatchAddOperationResponse) OperationId() string {
	return l.Data.OperationId
}

func (l *LatchStatusResponse) GetParentOperation() (operation LatchOperationStatus) {
	for _, operation = range l.Data.Operations {
		break
	}
	return
}

func (l *LatchStatusResponse) Status() string {
	return l.GetParentOperation().Status
}

func (l *LatchStatusResponse) TwoFactor() LatchTwoFactor {
	return l.GetParentOperation().TwoFactor
}

func (l *LatchStatusResponse) Operations() map[string]LatchOperationStatus {
	return l.GetParentOperation().Operations
}

func (l *LatchShowOperationResponse) Operations() (operations map[string]LatchOperation) {
	return l.Data.Operations
}

func (l *LatchShowOperationResponse) FirstOperation() (operationId string, operation LatchOperation) {
	for operationId, operation = range l.Data.Operations {
		break
	}
	return
}

func (l *LatchHistoryResponse) Application() LatchApplication {
	return l.Data.Application
}

func (l *LatchHistoryResponse) LastSeen() int64 {
	return l.Data.LastSeen
}

func (l *LatchHistoryResponse) ClientVersion() []LatchClientVersion {
	return l.Data.ClientVersion
}

func (l *LatchHistoryResponse) HistoryCount() int {
	return l.Data.HistoryCount
}

func (l *LatchHistoryResponse) History() []LatchHistoryEntry {
	return l.Data.History
}

func (l *LatchShowApplicationsResponse) Unmarshal(Json string) (err error) {
	return json.Unmarshal([]byte(Json), l)
}

func (l *LatchShowApplicationsResponse) Applications() map[string]LatchApplicationInfo {
	return l.Data.Applications
}

func (l *LatchAddApplicationResponse) Unmarshal(Json string) (err error) {
	return json.Unmarshal([]byte(Json), l)
}

func (l *LatchAddApplicationResponse) AppID() string {
	return l.Data.AppID
}

func (l *LatchAddApplicationResponse) Secret() string {
	return l.Data.Secret
}

func (l *LatchSubscriptionResponse) Unmarshal(Json string) (err error) {
	return json.Unmarshal([]byte(Json), l)
}

func (l *LatchSubscriptionResponse) ID() string {
	return l.Data.Subscription.ID
}

func (l *LatchSubscriptionResponse) Applications() LatchSubscriptionUsage {
	return l.Data.Subscription.Applications
}

func (l *LatchSubscriptionResponse) Operations() map[string]LatchSubscriptionUsage {
	return l.Data.Subscription.Operations
}

func (l *LatchSubscriptionResponse) Users() LatchSubscriptionUsage {
	return l.Data.Subscription.Users
}

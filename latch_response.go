package golatch

import (
	"encoding/json"
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

type LatchOperationResponse struct {
	Data struct {
		OperationId string `json:"operationId"`
	} `json:"data"`
}

type LatchOperationStatus struct {
	Status     string                          `json:"status"`
	TwoFactor  LatchTwoFactor                  `json:"two_factor"`
	Operations map[string]LatchOperationStatus `json:"operations"`
}

type LatchTwoFactor struct {
	Token     string `json:"token"`
	Generated string `json:"generated"`
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

func (l *LatchOperationResponse) Unmarshal(Json string) (err error) {
	return json.Unmarshal([]byte(Json), l)
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

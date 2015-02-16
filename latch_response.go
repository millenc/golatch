package golatch

import (
	"encoding/json"
)

type LatchResponse interface {
	Unmarshal(Json string) error
}

type LatchPairResponse struct {
	Data struct {
		AccountId string `json:"accountId"`
	} `json:"data"`
}

func (l *LatchPairResponse) Unmarshal(Json string) error {
	return json.Unmarshal([]byte(Json), l)
}

package golatch

import (
	"testing"
)

func TestLatchErrorResponseUnmarshal(t *testing.T) {
	json := `{"error":{"code":205, "message":"Account and application already paired"}}`
	response := &LatchErrorResponse{}
	err := response.Unmarshal(json)

	if err != nil {
		t.Errorf("LatchErrorResponse.Unmarshal() failed json: json: %q , error %q", json, err)
	} else if response.Err.Code != 205 || response.Err.Message != "Account and application already paired" {
		t.Errorf("LatchErrorResponse.Unmarshal() failed: json: %q , got %q", json, response)
	}
}

func TestLatchPairResponseUnmarshal(t *testing.T) {
	json := `{"data":{"accountId":"MyAccountId"}}`
	response := &LatchPairResponse{}
	err := response.Unmarshal(json)

	if err != nil {
		t.Errorf("LatchPairResponse.Unmarshal() failed json: json: %q , error %q", json, err)
	} else if response.Data.AccountId != "MyAccountId" {
		t.Errorf("LatchPairResponse.Unmarshal() failed: json: %q , got %q", json, response)
	}
}

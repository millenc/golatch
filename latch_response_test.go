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
	} else if response.AccountId() != "MyAccountId" {
		t.Errorf("LatchPairResponse.Unmarshal() failed: json: %q , got %q", json, response)
	}
}

func TestLatchStatusResponseUnmarshal(t *testing.T) {
	json := `{"data":{"operations":{"MyApplicationID":{"status":"on", "two_factor":{"token":"MyToken", "generated":"GeneratedTime"},"operations":{"MyOperationID":{"status":"on"}}}}}}`
	response := &LatchStatusResponse{}
	err := response.Unmarshal(json)

	if err != nil {
		t.Errorf("LatchStatusResponse.Unmarshal() failed json: json: %q , error %q", json, err)
	} else if response.Status() != LATCH_STATUS_ON {
		t.Errorf("LatchStatusResponse.Unmarshal() failed, expected on status: json: %q , got %q", json, response)
	} else if two_factor := response.TwoFactor(); two_factor.Token != "MyToken" || two_factor.Generated != "GeneratedTime" {
		t.Errorf("LatchStatusResponse.Unmarshal() failed, two factor data is wrong: json: %q , got %q", json, response)
	}

	operations := response.Operations()
	if operations == nil || len(operations) == 0 {
		t.Errorf("LatchStatusResponse.Unmarshal() failed, expected 1 operation, found none: json: %q , got %q", json, response)
		return
	}

	var key string
	var operation LatchOperationStatus
	for key, operation = range operations {
		break
	}
	if key != "MyOperationID" || operation.Status != LATCH_STATUS_ON {
		t.Errorf("LatchStatusResponse.Unmarshal() failed, expected 1 operation with ID %q and status %q: json: %q , got %q", "MyOperationID", "on", json, response)
		return
	}
}

func TestLatchAddOperationResponse(t *testing.T) {
	json := `{"data":{"operationId":"MyOperationId"}}`
	response := &LatchAddOperationResponse{}
	err := response.Unmarshal(json)

	if err != nil {
		t.Errorf("LatchAddOperationResponse.Unmarshal() failed json: %q , error %q", json, err)
	} else if response.OperationId() != "MyOperationId" {
		t.Errorf("LatchAddOperationResponse.Unmarshal() failed, expected operationId=%q : json: %q , got %q", "MyOperationId", json, response)
	}
}

func TestLatchShowOperationResponseUnmarshal(t *testing.T) {
	json := `{"data":{"operations":{"MyOperationId":{"name":"My Operation", "two_factor": "MANDATORY", "lock_on_request":"OPT_IN" ,"operations":{"MyNestedOperationID":{"name":"My Nested Operation"}}}}}}`
	response := &LatchShowOperationResponse{}
	err := response.Unmarshal(json)

	if err != nil {
		t.Errorf("LatchShowOperationResponse.Unmarshal() failed json: %q , error %q", json, err)
	}
	operation := response.Operation()
	if operation.Name != "My Operation" || operation.TwoFactor != MANDATORY || operation.LockOnRequest != OPT_IN || len(operation.Operations) == 0 {
		t.Errorf("LatchShowOperationResponse.Unmarshal() failed: expected:%q,%q,%q,%d and got %q,%q,%q,%d", "My Operation", MANDATORY, OPT_IN, 1, operation.Name, operation.TwoFactor, operation.LockOnRequest, len(operation.Operations))
	}

	var nested_operation_id string
	var nested_operation LatchOperation
	for nested_operation_id, nested_operation = range operation.Operations {
		break
	}

	if nested_operation_id != "MyNestedOperationID" || nested_operation.Name != "My Nested Operation" {
		t.Errorf("LatchShowOperationResponse.Unmarshal() failed: expected nested operation:%q with name %q, got %q with name %q", "MyNestedOperationID", "My Nested Operation", nested_operation_id, nested_operation.Name)
	}
}

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
	json := `{"data":{"operations":{"MyApplicationID":{"status":"on", "two_factor":{"token":"g2sEXg","generated":1425209705208},"operations":{"MyOperationID":{"status":"on"}}}}}}`
	response := &LatchStatusResponse{}
	err := response.Unmarshal(json)

	if err != nil {
		t.Errorf("LatchStatusResponse.Unmarshal() failed json: json: %q , error %q", json, err)
	} else if response.Status() != LATCH_STATUS_ON {
		t.Errorf("LatchStatusResponse.Unmarshal() failed, expected on status: json: %q , got %q", json, response)
	} else if two_factor := response.TwoFactor(); two_factor.Token != "g2sEXg" || two_factor.Generated != 1425209705208 {
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

func TestLatchHistoryResponseUnmarshal(t *testing.T) {
	json := `{"data":{"2Wv8UqaT6iZRQEbyG9Kv":{"status":"on","pairedOn":1428528090941,"name":"GoLatch Test","description":"","imageURL":"https://s3-eu-west-1.amazonaws.com/latch-ireland/avatar1.jpg","contactPhone":"666111222","contactEmail":"","two_factor":"DISABLED","lock_on_request":"DISABLED","operations":{"wJrfCBzZCtiZfVFwt9aJ":{"name":"Operation 1","status":"on","two_factor":"off","lock_on_request":"off","operations":{}}}},"lastSeen":1428858456785,"clientVersion":{"Android":"1.4.1"},"count":5,"history":[{"t":1428528254424,"action":"get","what":"status","value":"on","was":"-","name":"GoLatch Test","userAgent":"Go 1.1 package http","ip":"127.0.0.1"},{"t":1428528260264,"action":"USER_UPDATE","what":"status","value":"off","was":"on","name":"GoLatch Test","userAgent":"","ip":"127.0.0.1"},{"t":1428528264520,"action":"get","what":"status","value":"off","was":"-","name":"GoLatch Test","userAgent":"Go 1.1 package http","ip":"127.0.0.1"},{"t":1428528274326,"action":"USER_UPDATE","what":"status","value":"on","was":"off","name":"GoLatch Test","userAgent":"","ip":"127.0.0.1"},{"t":1428528277313,"action":"get","what":"status","value":"on","was":"-","name":"GoLatch Test","userAgent":"Go 1.1 package http","ip":"127.0.0.1"}]}}`
	response := &LatchHistoryResponse{AppID: "2Wv8UqaT6iZRQEbyG9Kv"}

	err := response.Unmarshal(json)

	if err != nil {
		t.Errorf("LatchHistoryResponse.Unmarshal() failed json: %q , error %q", json, err)
	}

	application := response.Application()
	operations := application.Operations
	lastSeen := response.LastSeen()
	clientVersion := response.ClientVersion()
	historyCount := response.HistoryCount()
	history := response.History()

	//Test application data
	if application.Status != "on" ||
		application.PairedOn != 1428528090941 ||
		application.Name != "GoLatch Test" ||
		application.Description != "" ||
		application.ImageURL != "https://s3-eu-west-1.amazonaws.com/latch-ireland/avatar1.jpg" ||
		application.ContactPhone != "666111222" ||
		application.ContactEmail != "" ||
		application.TwoFactor != DISABLED ||
		application.LockOnRequest != DISABLED {
		t.Errorf("LatchHistoryResponse.Unmarshal() failed, incorrect application data json: %s , object %s", json, response)
	}
	if operation := operations["wJrfCBzZCtiZfVFwt9aJ"]; len(operations) != 1 ||
		operation.Name != "Operation 1" ||
		operation.Status != "on" ||
		operation.LockOnRequest != "off" ||
		operation.TwoFactor != "off" {
		t.Errorf("LatchHistoryResponse.Unmarshal() failed, incorrect operations data json: %s , object %s", json, response)
	}

	//Test LastSeen
	if lastSeen != 1428858456785 {
		t.Errorf("LatchHistoryResponse.Unmarshal() failed, incorrect lastSeen data json: %s , object %s", json, response)
	}

	//Test Client Version
	if client := clientVersion["Android"]; len(clientVersion) != 1 || client != "1.4.1" {
		t.Errorf("LatchHistoryResponse.Unmarshal() failed, incorrect clientVersion data json: %s , object %s", json, response)
	}

	//Test History
	if historyCount != 5 || len(history) != 5 {
		t.Errorf("LatchHistoryResponse.Unmarshal() failed, incorrect history data json: %s , object %s", json, response)
	} else if firstHistoryEntry := history[0]; firstHistoryEntry.Time != 1428528254424 ||
		firstHistoryEntry.Action != "get" ||
		firstHistoryEntry.What != "status" ||
		firstHistoryEntry.Value != "on" ||
		firstHistoryEntry.Was != "-" ||
		firstHistoryEntry.Name != "GoLatch Test" ||
		firstHistoryEntry.UserAgent != "Go 1.1 package http" ||
		firstHistoryEntry.IP != "127.0.0.1" {
		t.Errorf("LatchHistoryResponse.Unmarshal() failed, incorrect history entry data json: %s , object %s", json, response)
	}
}

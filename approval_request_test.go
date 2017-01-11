package authy

import (
	"net/url"
	"testing"
	"time"
)

func Test_SendApprovalRequest(t *testing.T) {
	api := NewAuthyAPI(data.ApiKey)

	user, err := api.RegisterUser(data.Email, data.CountryCode, data.PhoneNumber, url.Values{})
	approvalRequest, err := api.SendApprovalRequest(user.ID, "please approve this", Details{"data1": "value1"}, url.Values{})

	if err != nil {
		t.Error("External error found", err)
	}

	if !approvalRequest.Valid() {
		t.Error("Apprval request should be valid.")
	}
}

func Test_FindApprovalRequest(t *testing.T) {
	api := NewAuthyAPI(data.ApiKey)

	user, err := api.RegisterUser(data.Email, data.CountryCode, data.PhoneNumber, url.Values{})
	approvalRequest, err := api.SendApprovalRequest(user.ID, "please approve this", Details{"data1": "value1"}, url.Values{})

	if err != nil {
		t.Error("External error found", err)
	}

	if !approvalRequest.Valid() {
		t.Error("Apprval request should be valid.")
	}

	uuid := approvalRequest.UUID
	approvalRequest, err = api.FindApprovalRequest(uuid, url.Values{})

	if err != nil {
		t.Error("External error found", err)
	}

	if approvalRequest.Status != "pending" {
		t.Error("Approval request status is wrong")
	}

	if uuid != approvalRequest.UUID {
		t.Error("Approval request doesn't match.")
	}
}

func Test_WaitForApprovalRequest(t *testing.T) {
	api := NewAuthyAPI(data.ApiKey)

	user, err := api.RegisterUser(data.Email, data.CountryCode, data.PhoneNumber, url.Values{})
	approvalRequest, err := api.SendApprovalRequest(user.ID, "please approve this", Details{"data1": "value1"}, url.Values{})

	if err != nil {
		t.Error("error found", err)
	}

	if !approvalRequest.Valid() {
		t.Error("Apprval request should be valid.")
	}

	now := time.Now()
	status, err := api.WaitForApprovalRequest(approvalRequest.UUID, 1*time.Second, url.Values{"user_ip": {"234.78.25.2"}})
	elapsedTime := time.Now().Sub(now)

	if err != nil {
		t.Error("error found", err)
	}

	if elapsedTime < 1 {
		t.Error("max duration not reached")
	}

	if status != OneTouchStatusExpired {
		t.Error("expired status expected")
	}
}

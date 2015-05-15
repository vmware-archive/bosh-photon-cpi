package esxcloud

import (
	"testing"
)

// Simple preliminary test. Make sure status API correctly deserializes the response
func TestGetStatus200(t *testing.T) {
	expectedStruct := Status{"READY", []Component{{"chairman", "", "READY"}, {"housekeeper", "", "READY"}}}
	server, client := testSetup()
	server.SetResponseJson(200, expectedStruct)
	defer server.Close()

	status, err := client.Status.Get()
	if err != nil {
		t.Error("Not expecting error from GetStatus")
	}

	if status.Status != expectedStruct.Status {
		t.Error("Status did not match expected result")
	}
	if len(status.Components) < 1 {
		t.Error("Expected to receive more than one status component")
		t.Log(status)
	}

	return
}

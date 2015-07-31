package esxcloud

import (
	"reflect"
	"testing"
)

func TestResourceTicketGetTask(t *testing.T) {
	if isIntegrationTest() {
		t.Skip("Skipping on integration mode. Need fix to run in integration")
	}
	mockTask := createMockTask("CREATE_RESOURCE_TICKET", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
	defer server.Close()

	taskList, err := client.ResourceTickets.GetTasks(mockTask.Entity.ID, nil)
	if err != nil {
		t.Error("Did not expect error from GetTasks")
		t.Log(err)
	}

	found := false
	for _,task := range taskList.Items {
		if reflect.DeepEqual(task, *mockTask) {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find task with resource ticket id " + mockTask.Entity.ID + " with state COMPLETED")
	}
}

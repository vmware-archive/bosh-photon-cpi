package esxcloud

import (
	"reflect"
	"testing"
)

func TestProjectGetTask(t *testing.T) {
	mockTask := createMockTask("CREATE_VM", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
	defer server.Close()

	taskList, err := client.Projects.GetTasks(mockTask.Entity.ID, nil)
	if err != nil {
		t.Error("Did not expect error from GetTasks")
	}

	found := false
	for _,task := range taskList.Items {
		if reflect.DeepEqual(task, *mockTask) {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find task with project id " + mockTask.Entity.ID + " with state COMPLETED")
	}
}

package esxcloud

import (
	"testing"
	"reflect"
)

func TestCreateGetAndDeleteHosts(t *testing.T) {
	if isIntegrationTest() {
		t.Skip("Skipping Host test on integration mode. No valid Host Address")
	}
	// Test Create
	mockTask := createMockTask("CREATE_HOST", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	hostSpec := &HostCreateSpec{
		Username: randomString(10),
		Password: randomString(10),
		Address: randomString(10),
		Tags: []string{},
	}
	createTask, err := client.Hosts.Create(hostSpec)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if createTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if createTask.Operation != "CREATE_HOST" {
		t.Error("Expected task operation to be CREATE_HOST")
	}
	if createTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Test Get
	server.SetResponseJson(200, Host{Username: hostSpec.Username})
	host, err := client.Hosts.Get(createTask.Entity.ID)
	if err != nil {
		t.Error("Did not expect error from Get")
		t.Log(err)
	}
	if host.Username != hostSpec.Username && host.Password != hostSpec.Password {
		t.Error("Host returned by Get did not match spec")
	}

	// Test GetAll
	server.SetResponseJson(200, &Hosts{[]Host{Host{Username: hostSpec.Username, Password: hostSpec.Password}}})
	hostList, err := client.Hosts.GetAll()
	if err != nil {
		t.Error("Did not expect error from GetAll")
		t.Log(err)
	}
	var found bool
	for _, h := range hostList.Items {
		if h.Username == hostSpec.Username && h.Password == hostSpec.Password {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find host with username " + hostSpec.Username + " and password " + hostSpec.Password)
	}

	// Test Delete
	mockTask = createMockTask("DELETE_HOST", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	deleteTask, err := client.Hosts.Delete(createTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if deleteTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if deleteTask.Operation != "DELETE_HOST" {
		t.Error("Expected task operation to be DELETE_HOST")
	}
	if deleteTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}
}

func TestHostGetTask(t *testing.T) {
	if isIntegrationTest() {
		t.Skip("Skipping Host test on integration mode. No valid Host Address")
	}
	mockTask := createMockTask("CREATE_HOST", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	hostSpec := &HostCreateSpec{
		Username: randomString(10),
		Password: randomString(10),
		Address: randomString(10),
		Tags: []string{},
	}
	createTask, err := client.Hosts.Create(hostSpec)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if createTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if createTask.Operation != "CREATE_HOST" {
		t.Error("Expected task operation to be CREATE_HOST")
	}
	if createTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	server.SetResponseJson(200, &TaskList{[]Task{*createTask}})
	taskList, err := client.Hosts.GetTasks(createTask.Entity.ID, nil)
	if err != nil {
		t.Error("Did not expect error from GetTasks")
		t.Log(err)
	}

	var found bool
	for _,task := range taskList.Items {
		if reflect.DeepEqual(task, *createTask) {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find task with host id " + createTask.Entity.ID)
	}

	taskList, err = client.Hosts.GetTasks(createTask.Entity.ID, &TaskGetOptions{State: "COMPLETED"})
	if err != nil {
		t.Error("Did not expect error from GetTasks")
		t.Log(err)
	}

	found = false
	for _,task := range taskList.Items {
		if reflect.DeepEqual(task, *createTask) {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find task with host id " + createTask.Entity.ID + " with state COMPLETED")
	}
}

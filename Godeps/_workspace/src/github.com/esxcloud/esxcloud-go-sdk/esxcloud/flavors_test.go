package esxcloud

import (
	"testing"
)

func TestCreateGetAndDelete(t *testing.T) {
	// Test Create
	mockTask := createMockTask("CREATE_FLAVOR", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	flavorSpec := &FlavorCreateSpec{
		Name: randomString(10),
		Kind: "vm",
		Cost: []QuotaLineItem{QuotaLineItem{"GB", 16, "vm.memory"}},
	}
	createTask, err := client.Flavors.Create(flavorSpec)
	if err != nil {
		t.Error("Not expecting error")
	}
	if createTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if createTask.Operation != "CREATE_FLAVOR" {
		t.Error("Expected task operation to be CREATE_FLAVOR")
	}
	if createTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Test Get
	server.SetResponseJson(200, Flavor{Name: flavorSpec.Name})
	flavor, err := client.Flavors.Get(createTask.Entity.ID)
	if err != nil {
		t.Error("Did not expect error from Get")
	}
	if flavor.Name != flavorSpec.Name {
		t.Error("Flavor returned by Get did not match spec")
	}

	// Test Get
	server.SetResponseJson(200, &FlavorList{[]Flavor{Flavor{Name: flavorSpec.Name}}})
	flavorList, err := client.Flavors.GetAll(&FlavorGetOptions{Name: flavorSpec.Name})
	if err != nil {
		t.Error("Did not expect error from Get")
	}
	var found bool
	for _, f := range flavorList.Items {
		if f.Name == flavorSpec.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find flavor with name " + flavorSpec.Name)
	}

	// Test Delete
	mockTask = createMockTask("DELETE_FLAVOR", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	deleteTask, err := client.Flavors.Delete(createTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if deleteTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if deleteTask.Operation != "DELETE_FLAVOR" {
		t.Error("Expected task operation to be DELETE_FLAVOR")
	}
	if deleteTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}
}

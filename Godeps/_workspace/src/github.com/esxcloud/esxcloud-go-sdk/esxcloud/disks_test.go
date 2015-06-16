package esxcloud

import (
	"testing"
)

func TestCreateGetDeleteDisk(t *testing.T) {
	// Create tenant
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: randomString(10)}
	tenantTask, _ := client.Tenants.Create(tenantSpec)

	// Create resource ticket
	resSpec := &ResourceTicketCreateSpec{
		Name:   randomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 16, Key: "vm.memory"}},
	}
	_, _ = client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, resSpec)

	// Create project
	projSpec := &ProjectCreateSpec{
		ResourceTicketReservation{
			resSpec.Name,
			[]QuotaLineItem{QuotaLineItem{"GB", 2, "vm.memory"}},
		},
		randomString(10),
	}
	projTask, _ := client.Tenants.CreateProject(tenantTask.Entity.ID, projSpec)

	// Create flavor
	flavorSpec := &FlavorCreateSpec{
		[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "persistent-disk.cost"}},
		"persistent-disk",
		randomString(10),
	}
	_, _ = client.Flavors.Create(flavorSpec)

	// Create disk
	mockTask = createMockTask("CREATE_DISK", "QUEUED")
	server.SetResponseJson(200, mockTask)
	diskSpec := &DiskCreateSpec{
		Flavor:     flavorSpec.Name,
		Kind:       "persistent-disk",
		CapacityGB: 2,
		Name:       randomString(10),
	}
	diskTask, err := client.Projects.CreateDisk(projTask.Entity.ID, diskSpec)
	if err != nil {
		t.Error("Not expecting error")
	}
	if diskTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if diskTask.Operation != "CREATE_DISK" {
		t.Error("Expected task operation to be CREATE_DISK")
	}
	if diskTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for disk creation to complete
	mockTask = createMockTask("CREATE_DISK", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	diskTask, err = client.Tasks.Wait(diskTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if diskTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if diskTask.Operation != "CREATE_DISK" {
		t.Error("Expected task operation to be CREATE_DISK")
	}
	if diskTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Get disk
	diskMock := &PersistentDisk{
		Name:       diskSpec.Name,
		Flavor:     diskSpec.Flavor,
		CapacityGB: diskSpec.CapacityGB,
		Kind:       diskSpec.Kind,
	}
	server.SetResponseJson(200, diskMock)
	disk, err := client.Disks.Get(diskTask.Entity.ID)
	if disk.Flavor != diskSpec.Flavor {
		t.Error("Disk flavor did not match spec")
	}
	if disk.Name != diskSpec.Name {
		t.Error("Disk name did not match spec")
	}
	if disk.Kind != diskSpec.Kind {
		t.Error("Disk kind did not match spec")
	}
	if disk.CapacityGB != diskSpec.CapacityGB {
		t.Error("Disk capacity did not match spec")
	}

	// Delete disk
	mockTask = createMockTask("DELETE_DISK", "QUEUED")
	server.SetResponseJson(200, mockTask)
	deleteTask, err := client.Disks.Delete(diskTask.Entity.ID, false)
	if err != nil {
		t.Error("Not expecting error")
	}
	if deleteTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if deleteTask.Operation != "DELETE_DISK" {
		t.Error("Expected task operation to be DELETE_DISK")
	}
	if deleteTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for disk deletion
	mockTask = createMockTask("DELETE_DISK", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	deleteTask, err = client.Tasks.Wait(deleteTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if deleteTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if deleteTask.Operation != "DELETE_DISK" {
		t.Error("Expected task operation to be DELETE_DISK")
	}
	if deleteTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Cleanup project
	_, err = client.Projects.Delete(projTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting project")
		t.Log(err)
	}

	// Cleanup tenant
	_, err = client.Tenants.Delete(tenantTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting tenant")
		t.Log(err)
	}
}

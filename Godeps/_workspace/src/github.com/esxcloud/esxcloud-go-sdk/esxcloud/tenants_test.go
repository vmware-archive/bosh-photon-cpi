package esxcloud

import (
	"reflect"
	"testing"
)

func TestCreateGetAndDeleteTenants(t *testing.T) {
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: randomString(10, "go-sdk-tenant-")}
	task, err := client.Tenants.Create(tenantSpec)
	if err != nil {
		t.Error("Not expecting error from Create")
		t.Log(err)
	}
	if task == nil {
		t.Error("Not expecting task to be nil")
	}
	if task.Operation != "CREATE_TENANT" {
		t.Error("Expected task operation to be CREATE_TENANT")
	}
	if task.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	server.SetResponseJson(200, &Tenants{Items: []Tenant{Tenant{Name: tenantSpec.Name}}})
	tenants, err := client.Tenants.GetAll()
	if err != nil {
		t.Error("Not expecting error from GetAll")
		t.Log(err)
	}
	if tenants == nil {
		t.Error("Not expecting tenants to be nil")
	}

	var found bool
	for _, tenant := range tenants.Items {
		if tenant.Name == tenantSpec.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find newly created tenant in result of GetAll")
	}

	mockTask = createMockTask("DELETE_TENANT", "COMPLETED")
	server.SetResponseJson(200, mockTask)

	task, err = client.Tenants.Delete(task.Entity.ID)
	if err != nil {
		t.Error("Not expecting error from Create")
		t.Log(err)
	}
	if task == nil {
		t.Error("Not expecting task to be nil")
	}
	if task.Operation != "DELETE_TENANT" {
		t.Error("Expected task operation to be DELETE_TENANT")
	}
	if task.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}
}

func TestCreateAndGetResTickets(t *testing.T) {
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: randomString(10, "go-sdk-tenant-")}
	tenantTask, _ := client.Tenants.Create(tenantSpec)
	spec := &ResourceTicketCreateSpec{
		Name:   randomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 16, Key: "vm.memory"}},
	}

	mockTask = createMockTask("CREATE_RESOURCE_TICKET", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	task, err := client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, spec)
	if err != nil {
		t.Error("Not expecting error from CreateResourceTicket")
		t.Log(err)
	}
	if task == nil {
		t.Error("Not expecting task to be nil")
	}
	if task.Operation != "CREATE_RESOURCE_TICKET" {
		t.Error("Expected task operation to be CREATE_RESOURCE_TICKET")
	}
	if task.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	mockResList := ResourceList{[]ResourceTicket{ResourceTicket{TenantId: tenantTask.Entity.ID, Name: spec.Name, Limits: spec.Limits}}}
	server.SetResponseJson(200, mockResList)
	resList, err := client.Tenants.GetResourceTickets(tenantTask.Entity.ID, &ResourceTicketGetOptions{spec.Name})
	if err != nil {
		t.Error("Not expecting error from GetResourceTickets")
		t.Log(err)
	}
	if resList == nil {
		t.Error("Not expecting resource list to be nil")
	}
	var found bool
	for _, res := range resList.Items {
		if res.Name == spec.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Resource ticket not found")
	}

	// Cleanup tenant
	_, err = client.Tenants.Delete(tenantTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting tenant")
		t.Log(err)
	}
}

func TestCreateAndGetProjects(t *testing.T) {
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: randomString(10, "go-sdk-tenant-")}
	tenantTask, _ := client.Tenants.Create(tenantSpec)

	resSpec := &ResourceTicketCreateSpec{
		Name:   randomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 16, Key: "vm.memory"}},
	}
	_, _ = client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, resSpec)

	projSpec := &ProjectCreateSpec{
		ResourceTicketReservation{
			resSpec.Name,
			[]QuotaLineItem{QuotaLineItem{"GB", 2, "vm.memory"}},
		},
		randomString(10, "go-sdk-project-"),
	}

	mockTask = createMockTask("CREATE_PROJECT", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	task, err := client.Tenants.CreateProject(tenantTask.Entity.ID, projSpec)

	if err != nil {
		t.Error("Not expecting error from CreateProject")
		t.Log(err)
	}
	if task == nil {
		t.Error("Not expecting task to be nil")
	}
	if task.Operation != "CREATE_PROJECT" {
		t.Error("Expected task operation to be CREATE_PROJECT")
	}
	if task.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	mockProjects := ProjectList{[]ProjectCompact{ProjectCompact{Name: projSpec.Name}}}
	server.SetResponseJson(200, mockProjects)
	projList, err := client.Tenants.GetProjects(tenantTask.Entity.ID, &ProjectGetOptions{projSpec.Name})
	if err != nil {
		t.Error("Not expecting error from GetProjects")
		t.Log(err)
	}
	if projList == nil {
		t.Error("Not expecting project list to be nil")
	}
	var found bool
	for _, proj := range projList.Items {
		if proj.Name == projSpec.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Project not found")
	}

	// Cleanup project
	_, err = client.Projects.Delete(task.Entity.ID)
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

func TestTenantGetTask(t *testing.T) {
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{
		Name: randomString(10, "go-sdk-tenant-"),
	}
	createTask, err := client.Tenants.Create(tenantSpec)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if createTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if createTask.Operation != "CREATE_TENANT" {
		t.Error("Expected task operation to be CREATE_TENANT")
	}
	if createTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	server.SetResponseJson(200, &TaskList{[]Task{*createTask}})
	taskList, err := client.Tenants.GetTasks(createTask.Entity.ID, nil)
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
		t.Error("Did not find task with tenant id " + createTask.Entity.ID)
	}

	taskList, err = client.Tenants.GetTasks(createTask.Entity.ID, &TaskGetOptions{State: "COMPLETED"})
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
		t.Error("Did not find task with tenant id " + createTask.Entity.ID + " with state COMPLETED")
	}
}

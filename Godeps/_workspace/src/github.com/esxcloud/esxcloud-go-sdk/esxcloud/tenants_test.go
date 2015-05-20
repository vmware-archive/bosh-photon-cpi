package esxcloud

import (
	"testing"
)

func TestCreateGetAndDeleteTenants(t *testing.T) {
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: RandomString(10)}
	task, err := client.Tenants.Create(tenantSpec)
	if err != nil {
		t.Error("Not expecting error from Create")
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

	tenantSpec := &TenantCreateSpec{Name: RandomString(10)}
	tenantTask, _ := client.Tenants.Create(tenantSpec)
	spec := &ResourceTicketCreateSpec{
		Name:   RandomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 16, Key: "vm.memory"}},
	}

	mockTask = createMockTask("CREATE_RESOURCE_TICKET", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	task, err := client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, spec)
	if err != nil {
		t.Error("Not expecting error from CreateResourceTicket")
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
	resList, err := client.Tenants.GetResourceTickets(tenantTask.Entity.ID, &spec.Name)
	if err != nil {
		t.Error("Not expecting error from GetResourceTickets")
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
	}
}

func TestCreateAndGetProjects(t *testing.T) {
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: RandomString(10)}
	tenantTask, _ := client.Tenants.Create(tenantSpec)

	resSpec := &ResourceTicketCreateSpec{
		Name:   RandomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 16, Key: "vm.memory"}},
	}
	_, _ = client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, resSpec)

	projSpec := &ProjectCreateSpec{
		ResourceTicketReservation{
			resSpec.Name,
			[]QuotaLineItem{QuotaLineItem{"GB", 2, "vm.memory"}},
		},
		RandomString(10),
	}

	mockTask = createMockTask("CREATE_PROJECT", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	task, err := client.Tenants.CreateProject(tenantTask.Entity.ID, projSpec)

	if err != nil {
		t.Error("Not expecting error from CreateProject")
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
	projList, err := client.Tenants.GetProjects(tenantTask.Entity.ID, &projSpec.Name)
	if err != nil {
		t.Error("Not expecting error from GetProjects")
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
	}

	// Cleanup tenant
	_, err = client.Tenants.Delete(tenantTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting tenant")
	}
}
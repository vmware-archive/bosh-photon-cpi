package esxcloud

import (
	"testing"
)

func TestCreateGetAndDeleteDeployment(t *testing.T) {
	if isIntegrationTest() {
		t.Skip("Skipping deployment test on integration mode. Need undeployed environment")
	}
	// Test Create
	mockTask := createMockTask("CREATE_DEPLOYMENT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	deploymentSpec := &DeploymentCreateSpec{
		ImageDatastore: randomString(10),
		Auth: &AuthInfo {
			Enabled: false,
		},
	}
	createTask, err := client.Deployments.Create(deploymentSpec)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if createTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if createTask.Operation != "CREATE_DEPLOYMENT" {
		t.Error("Expected task operation to be CREATE_DEPLOYMENT")
	}
	if createTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Test Get
	server.SetResponseJson(200, Deployment{ImageDatastore: deploymentSpec.ImageDatastore, Auth: &AuthInfo{Enabled: false}})
	deployment, err := client.Deployments.Get(createTask.Entity.ID)
	if err != nil {
		t.Error("Did not expect error from Get")
		t.Log(err)
	}
	if deployment.ImageDatastore != deploymentSpec.ImageDatastore {
		t.Error("Deployment returned by Get did not match spec")
	}

	// Test GetAll
	server.SetResponseJson(200, &Deployments{[]Deployment{Deployment{ImageDatastore: deploymentSpec.ImageDatastore, Auth: &AuthInfo{Enabled: false}}}})
	deploymentList, err := client.Deployments.GetAll()
	if err != nil {
		t.Error("Did not expect error from GetAll")
		t.Log(err)
	}
	var found bool
	for _, d := range deploymentList.Items {
		if d.ImageDatastore == deploymentSpec.ImageDatastore {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find deployment with image datastore name " + deploymentSpec.ImageDatastore)
	}

	// Test Delete
	mockTask = createMockTask("DELETE_DEPLOYMENT", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	deleteTask, err := client.Deployments.Delete(createTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if deleteTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if deleteTask.Operation != "DELETE_DEPLOYMENT" {
		t.Error("Expected task operation to be DELETE_DEPLOYMENT")
	}
	if deleteTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}
}

func TestDeploymentGetHosts(t *testing.T) {
	if isIntegrationTest() {
		t.Skip("Skipping deployment test on integration mode. Need undeployed environment")
	}
	mockTask := createMockTask("CREATE_DEPLOYMENT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	deploymentSpec := &DeploymentCreateSpec{
		ImageDatastore: randomString(10),
		Auth: &AuthInfo {
			Enabled: false,
		},
	}
	createTask, err := client.Deployments.Create(deploymentSpec)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if createTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if createTask.Operation != "CREATE_DEPLOYMENT" {
		t.Error("Expected task operation to be CREATE_DEPLOYMENT")
	}
	if createTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	mockTask = createMockTask("CREATE_HOST", "COMPLETED")
	server.SetResponseJson(200, mockTask)

	hostSpec := &HostCreateSpec{
		Username: randomString(10),
		Password: randomString(10),
		Address: randomString(10),
		Tags: []string{},
	}

	server.SetResponseJson(200, &Hosts{[]Host{Host{Username: hostSpec.Username, Password: hostSpec.Password}}})
	hostList, err := client.Deployments.GetHosts(createTask.Entity.ID)
	if err != nil {
		t.Error("Did not expect error from GetHosts")
		t.Log(err)
	}
	var found bool
	for _, host := range hostList.Items {
		if host.Username == hostSpec.Username && host.Password == hostSpec.Password {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find host with deployment id " + createTask.Entity.ID + " host name " + hostSpec.Username + " host password " + hostSpec.Password)
	}
}

func TestDeploymentGetVMs(t *testing.T) {
	if isIntegrationTest() {
		t.Skip("Skipping deployment test on integration mode. Need undeployed environment")
	}
	mockTask := createMockTask("CREATE_DEPLOYMENT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	deploymentSpec := &DeploymentCreateSpec{
		ImageDatastore: randomString(10),
		Auth: &AuthInfo {
			Enabled: false,
		},
	}
	createTask, err := client.Deployments.Create(deploymentSpec)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if createTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if createTask.Operation != "CREATE_DEPLOYMENT" {
		t.Error("Expected task operation to be CREATE_DEPLOYMENT")
	}
	if createTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	mockTask = createMockTask("CREATE_VM", "COMPLETED")
	server.SetResponseJson(200, mockTask)

	vmSpec := &VM{
		Name: randomString(10, "go-sdk-vm-"),
	}

	server.SetResponseJson(200, &VMs{[]VM{VM{Name: vmSpec.Name}}})
	vmList, err := client.Deployments.GetVms(createTask.Entity.ID)
	if err != nil {
		t.Error("Did not expect error from GetHosts")
		t.Log(err)
	}
	var found bool
	for _, vm := range vmList.Items {
		if vm.Name == vmSpec.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find vm with deployment id " + createTask.Entity.ID + " vm name " + vmSpec.Name)
	}
}

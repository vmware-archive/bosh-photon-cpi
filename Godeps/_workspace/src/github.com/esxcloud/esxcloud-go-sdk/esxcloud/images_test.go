package esxcloud

import (
	"testing"
)

func TestCreateAndDeleteImage(t *testing.T) {
	mockTask := createMockTask("CREATE_IMAGE", "STARTED", createMockStep("UPLOAD_IMAGE", "COMPLETED"))
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	// Start create image
	imagePath := "../testdata/tty_tiny.ova"
	task, err := client.Images.CreateFromFile(imagePath, nil)
	if err != nil {
		t.Error("Not expecting error from image create")
	}
	if task == nil {
		t.Error("Not expecting task to be nil")
	}
	if task.Operation != "CREATE_IMAGE" {
		t.Error("Expected task operation to be CREATE_IMAGE")
	}
	if task.State != "STARTED" {
		t.Error("Expected task status to be STARTED")
	}
	if !hasStep(task, "UPLOAD_IMAGE", "COMPLETED") {
		t.Error("Expected to find a task UPLOAD_IMAGE with state COMPLETED")
	}

	// Wait for create image to be completed
	mockTask = createMockTask("CREATE_IMAGE", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	task, err = client.Tasks.Wait(task.ID)
	if err != nil {
		t.Error("Not expecting error from image create")
	}
	if task == nil {
		t.Error("Not expecting task to be nil")
	}
	if task.Operation != "CREATE_IMAGE" {
		t.Error("Expected task operation to be CREATE_IMAGE")
	}
	if task.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Start delete image
	mockTask = createMockTask("DELETE_IMAGE", "QUEUED")
	server.SetResponseJson(200, mockTask)
	task, err = client.Images.Delete(task.Entity.ID)
	if err != nil {
		t.Error("Not expecting error from image delete")
	}
	if task == nil {
		t.Error("Not expecting task to be nil")
	}
	if task.Operation != "DELETE_IMAGE" {
		t.Error("Expected task operation to be DELETE_IMAGE")
	}
	if task.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for delete image to be completed
	mockTask = createMockTask("DELETE_IMAGE", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	task, err = client.Tasks.Wait(task.ID)
	if err != nil {
		t.Error("Not expecting error from image delete")
	}
	if task == nil {
		t.Error("Not expecting task to be nil")
	}
	if task.Operation != "DELETE_IMAGE" {
		t.Error("Expected task operation to be DELETE_IMAGE")
	}
	if task.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}
}

package esxcloud

import (
	"testing"
	"time"
)

func TestTaskError(t *testing.T) {
	// Unit test only
	if isIntegrationTest() {
		return
	}
	server, client := testSetup()
	defer server.Close()

	task := &Task{ID: "fake-id", State: "ERROR", Operation: "fake-op"}
	server.SetResponseJson(200, task)
	task, err := client.Tasks.Wait(task.ID)
	taskErr, ok := err.(TaskError)
	if !ok {
		t.Error("Expecting to find error of type TaskError")
	} else if taskErr.ID != task.ID {
		t.Error("TaskError.ID not equal to actual task ID")
	}

	client.options.TaskPollTimeout = 1 * time.Second
}

func TestTaskTimeoutError(t *testing.T) {
	// Unit test only
	if isIntegrationTest() {
		return
	}
	server, client := testSetup()
	defer server.Close()

	client.options.TaskPollTimeout = 1 * time.Second
	task := &Task{ID: "fake-id", State: "QUEUED", Operation: "fake-op"}
	server.SetResponseJson(200, task)
	task, err := client.Tasks.Wait(task.ID)
	taskErr, ok := err.(TaskTimeoutError)
	if !ok {
		t.Error("Expecting to find error of type TaskTimeoutError")
	} else if taskErr.ID != task.ID {
		t.Error("TaskTimeoutError.ID not equal to actual task ID")
	}
}

func TestHttpError(t *testing.T) {
	// Unit test only
	if isIntegrationTest() {
		return
	}
	server, client := testSetup()
	defer server.Close()

	client.options.TaskPollTimeout = 1 * time.Second
	task := &Task{ID: "fake-id", State: "QUEUED", Operation: "fake-op"}
	server.SetResponseJson(500, "server error")
	task, err := client.Tasks.Wait(task.ID)
	taskErr, ok := err.(HttpError)
	if !ok {
		t.Error("Expecting to find error of type HttpError")
	} else if taskErr.StatusCode != 500 {
		t.Error("HttpError.StatusCode did not equal error from server")
	}
}

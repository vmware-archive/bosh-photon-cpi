package esxcloud

import (
	"encoding/json"
	"time"

	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

// Contains functionality for tasks API.
type TasksAPI struct {
	client *Client
}

var TaskUrl string = "/tasks"

// Gets a task by ID.
func (api *TasksAPI) Get(id string) (task *Task, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+TaskUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := getTask(getError(res))
	return result, err
}

// Gets all tasks, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *TasksAPI) GetAll(options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint+TaskUrl
	if options != nil {
		uri += getQueryString(options)
	}
	res, err := rest.Get(api.client.httpClient, uri, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &TaskList{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Waits for a task to complete by polling the tasks API until a task returns with
// the state COMPLETED or ERROR. Will wait no longer than the duration specified by timeout.
func (api *TasksAPI) WaitTimeout(id string, timeout time.Duration) (task *Task, err error) {
	start := time.Now()
	numErrors := 0
	maxErrors := api.client.options.taskRetryCount

	for time.Since(start) < timeout {
		task, err = api.Get(id)
		if err != nil {
			switch err.(type) {
			// If an ApiError comes back, something is wrong, return the error to the caller
			case ApiError:
				return
			// For other errors, retry before giving up
			default:
				numErrors++
				if numErrors > maxErrors {
					return
				}
			}
		} else {
			// Reset the error count any time a successful call is made
			numErrors = 0
			if task.State == "COMPLETED" {
				return
			}
			if task.State == "ERROR" {
				err = TaskError{id}
				return
			}
		}
		time.Sleep(api.client.options.taskPollDelay)
	}
	err = TaskTimeoutError{id}
	return
}

// Waits for a task to complete by polling the tasks API until a task returns with
// the state COMPLETED or ERROR.
func (api *TasksAPI) Wait(id string) (task *Task, err error) {
	return api.WaitTimeout(id, api.client.options.TaskPollTimeout)
}

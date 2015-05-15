package esxcloud

import (
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
	"math"
	"time"
)

type TasksAPI struct {
	client *Client
}

func (api *TasksAPI) Get(id string) (task *Task, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/tasks/"+id)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := getTask(getError(res))
	return result, err
}

// Waits for a task to complete by polling the tasks API until a task returns with
// the state COMPLETED or ERROR. Will wait no longer than the duration specified by timeout.
func (api *TasksAPI) WaitTimeout(id string, timeout time.Duration) (task *Task, err error) {
	start := time.Now()
	numErrors := 0
	maxErrors := 3

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
		}
		// Reset the error count any time a successful call is made
		numErrors = 0
		if task.State == "COMPLETED" || task.State == "ERROR" {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

// Waits for a task to complete by polling the tasks API until a task returns with
// the state COMPLETED or ERROR.
func (api *TasksAPI) Wait(id string) (task *Task, err error) {
	return api.WaitTimeout(id, math.MaxInt64)
}

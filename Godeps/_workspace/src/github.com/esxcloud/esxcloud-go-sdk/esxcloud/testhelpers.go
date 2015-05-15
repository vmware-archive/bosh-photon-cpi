package esxcloud

import (
	"encoding/json"
)

func toJson(v interface{}) string {
	res, err := json.Marshal(v)
	if err != nil {
		// Since this method is only for testing, don't return
		// any errors, just panic.
		panic("Error serializing struct into JSON")
	}
	// json.Marshal returns []byte, convert to string
	return string(res[:])
}

func hasStep(task *Task, operation, state string) bool {
	for _, step := range task.Steps {
		if step.State == state && step.Operation == operation {
			return true
		}
	}
	return false
}

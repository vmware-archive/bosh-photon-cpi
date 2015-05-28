package esxcloud

import (
	"fmt"
)

type Entity struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

type ApiError struct {
	Code           string                 `json:"code"`
	Data           map[string]interface{} `json:"data"`
	Message        string                 `json:"message"`
	HttpStatusCode int                    `json:"-"` // Not part of API contract
}

// Implement Go error interface for ApiError
func (e ApiError) Error() string {
	return fmt.Sprintf(
		"esxcloud: { HTTP status: '%v', code: '%v', message: '%v', data: '%v' }",
		e.HttpStatusCode,
		e.Code,
		e.Message,
		e.Data)
}

// Used to represent a generic HTTP error
type HttpError struct {
	StatusCode int
	Message    string
}

// Implementation of error interface for HttpError
func (e HttpError) Error() string {
	return fmt.Sprintf("esxcloud: HTTP %d: %v", e.StatusCode, e.Message)
}

// Represents an ESXCloud task that has entered into an error state.
// ESXCloud task errors can be caught and type-checked against with
// the usual Go idiom.
type TaskError struct {
	ID string
}

func (e TaskError) Error() string {
	return fmt.Sprintf("esxcloud: Task '%s' is in error state. " +
		"Examine task for full details.")
}

// An error representing a timeout while waiting for a task to complete.
type TaskTimeoutError struct {
	ID string
}

func (e TaskTimeoutError) Error() string {
	return fmt.Sprintf("esxcloud: Timed out waiting for task '%s'. " +
		"Task may not be in error state, examine task for full details.")
}

type Step struct {
	ID                 string                 `json:"id"`
	Operation          string                 `json:"operation,omitempty"`
	State              string                 `json:"state"`
	StartedTime        int64                  `json:"startedTime"`
	EndTime            int64                  `json:"endTime,omitempty"`
	QueuedTime         int64                  `json:"queuedTime"`
	Sequence           int                    `json:"sequence,omitempty"`
	ResourceProperties map[string]interface{} `json:"resourceProperties,omitempty"`
	Errors             []ApiError             `json:"errors,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
	SelfLink           string                 `json:"selfLink,omitempty"`
}

type Task struct {
	ID          string `json:"id"`
	Operation   string `json:"operation,omitempty"`
	State       string `json:"state"`
	StartedTime int64  `json:"startedTime"`
	EndTime     int64  `json:"endTime,omitempty"`
	QueuedTime  int64  `json:"queuedTime"`
	Entity      Entity `json:"entity,omitempty"`
	SelfLink    string `json:"selfLink,omitempty"`
	Steps       []Step `json:"steps,omitempty"`
}

type BaseCompact struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type QuotaLineItem struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
	Key   string  `json:"key"`
}

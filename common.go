package main

import (
	. "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
)

// Indicates whether or not an error is of type esxcloud.TaskError
func isTaskError(e error) bool {
	if _, ok := e.(TaskError); ok {
		return true
	}
	return false
}

package main

import (
	. "github.com/esxcloud/photon-go-sdk/photon"
)

// Indicates whether or not an error is of type photon.TaskError
func isTaskError(e error) bool {
	if _, ok := e.(TaskError); ok {
		return true
	}
	return false
}

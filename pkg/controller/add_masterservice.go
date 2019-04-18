package controller

import (
	"github.com/blablacar/masterservice-operator/pkg/controller/masterservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, masterservice.Add)
}

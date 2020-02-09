package controller

import (
	"github.com/apirator/apirator/pkg/controller/apimock"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, apimock.Add)
}

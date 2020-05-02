package resources

import (
	apirator "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var predicateLogger = logf.Log.WithName("cmd")

type StatusChangedPredicate struct {
}

func (StatusChangedPredicate) Create(event.CreateEvent) bool {
	predicateLogger.Info("Creation event")
	return true
}

func (StatusChangedPredicate) Delete(event.DeleteEvent) bool {
	predicateLogger.Info("Deletion event")
	return true
}

func (StatusChangedPredicate) Generic(event.GenericEvent) bool {
	predicateLogger.Info("Generic event")
	return true
}

func (StatusChangedPredicate) Update(e event.UpdateEvent) bool {
	o, oldMock := e.ObjectOld.(*apirator.APIMock)
	n, newMock := e.ObjectNew.(*apirator.APIMock)
	if oldMock && newMock {
		return o.CountSteps() == n.CountSteps()
	}
	predicateLogger.Info("Update event")
	return true
}

type CreatePredicate struct {
}

func (CreatePredicate) Create(event.CreateEvent) bool {
	predicateLogger.Info("Create event")
	return true
}

func (CreatePredicate) Delete(event.DeleteEvent) bool {
	predicateLogger.Info("ignoring Delete event")
	return false
}

func (CreatePredicate) Update(event.UpdateEvent) bool {
	predicateLogger.Info("ignoring Update event")
	return false
}

func (CreatePredicate) Generic(event.GenericEvent) bool {
	predicateLogger.Info("ignoring Generic event")
	return false
}

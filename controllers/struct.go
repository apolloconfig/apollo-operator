package controllers

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReconcilerParams is the set of options to build a new Reconciler.
type ReconcilerParams struct { // 感觉可以改个名字，不叫 Params
	client.Client
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Tasks    []Task
}

// Task represents a reconciliation task to be executed by the reconciler.
type Task struct {
	Do          func(context.Context, client.Object, models.Params) error
	Name        string
	BailOnError bool
}

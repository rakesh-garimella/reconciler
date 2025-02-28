package example

import (
	"github.com/kyma-incubator/reconciler/pkg/reconciler"
	"github.com/kyma-incubator/reconciler/pkg/reconciler/service"
)

type CustomAction struct {
	name string
}

func (a *CustomAction) Run(version, profile string, config []reconciler.Configuration, context *service.ActionContext) error {
	if _, err := context.KubeClient.Clientset(); err != nil { //example how to retrieve native Kubernetes GO client
		context.Logger.Errorf("Failed to retrieve native Kubernetes GO client")
	}

	context.Logger.Infof("Action '%s' executed (passed version was '%s')", a.name, version)

	return nil
}

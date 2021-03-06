/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/trait"
	"github.com/apache/camel-k/pkg/util/kubernetes"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
)

// NewDeployAction create an action that handles integration deploy
func NewDeployAction() Action {
	return &deployAction{}
}

type deployAction struct {
}

func (action *deployAction) Name() string {
	return "deploy"
}

func (action *deployAction) CanHandle(integration *v1alpha1.Integration) bool {
	return integration.Status.Phase == v1alpha1.IntegrationPhaseDeploying
}

func (action *deployAction) Handle(integration *v1alpha1.Integration) error {
	resources, err := trait.ComputeDeployment(integration)
	if err != nil {
		return err
	}
	// TODO we should look for objects that are no longer present in the collection and remove them
	err = kubernetes.ReplaceResources(resources)
	if err != nil {
		return err
	}

	target := integration.DeepCopy()
	logrus.Info("Integration ", target.Name, " transitioning to state ", v1alpha1.IntegrationPhaseRunning)
	target.Status.Phase = v1alpha1.IntegrationPhaseRunning

	return sdk.Update(target)
}

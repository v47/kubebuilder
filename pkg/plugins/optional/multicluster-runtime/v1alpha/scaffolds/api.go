/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scaffolds

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins/optional/multicluster-runtime/v1alpha/scaffolds/internal/templates/controllers"
)

var _ plugins.Scaffolder = &apiScaffolder{}

type apiScaffolder struct {
	config   config.Config
	resource resource.Resource
	fs       machinery.Filesystem
}

// NewAPIScaffolder returns a Scaffolder for the create api command.
func NewAPIScaffolder(cfg config.Config, res resource.Resource) plugins.Scaffolder {
	return &apiScaffolder{config: cfg, resource: res}
}

// InjectFS implements plugins.Scaffolder.
func (s *apiScaffolder) InjectFS(fs machinery.Filesystem) { s.fs = fs }

// Scaffold overwrites the controller file with a multicluster-aware version.
func (s *apiScaffolder) Scaffold() error {
	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
		machinery.WithResource(&s.resource),
	)

	return scaffold.Execute(
		&controllers.Controller{Force: true},
	)
}

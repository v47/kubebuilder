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

package v1alpha

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins/optional/multicluster-runtime/v1alpha/scaffolds"
)

var _ plugin.CreateAPISubcommand = &createAPISubcommand{}

type createAPISubcommand struct {
	config   config.Config
	resource *resource.Resource
}

func (p *createAPISubcommand) UpdateMetadata(_ plugin.CLIMetadata, subcmdMeta *plugin.SubcommandMetadata) {
	subcmdMeta.Description = `Overwrite the controller scaffolded by go/v4 with a multicluster-aware version.

The generated controller uses:
  - mcreconcile.Request  (carries ClusterName alongside NamespacedName)
  - mcbuilder.ControllerManagedBy(mgr)  (watches across all registered clusters)
  - mcmanager.Manager  (multicluster manager type)`
	subcmdMeta.Examples = `  kubebuilder create api \
    --plugins go/v4,multicluster-runtime.sigs.k8s.io/v1 \
    --group foo --version v1 --kind Foo --controller --resource`
}

func (p *createAPISubcommand) InjectConfig(c config.Config) error {
	p.config = c
	return nil
}

func (p *createAPISubcommand) InjectResource(res *resource.Resource) error {
	p.resource = res
	return nil
}

func (p *createAPISubcommand) Scaffold(fs machinery.Filesystem) error {
	if p.resource == nil || !p.resource.HasController() {
		return nil
	}
	s := scaffolds.NewAPIScaffolder(p.config, *p.resource)
	s.InjectFS(fs)
	return s.Scaffold()
}

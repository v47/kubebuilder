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
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins/optional/multicluster-runtime/v1alpha/scaffolds"
)

var _ plugin.EditSubcommand = &editSubcommand{}

type editSubcommand struct {
	config        config.Config
	provider      string
	kubeconfigDir string
}

func (p *editSubcommand) UpdateMetadata(_ plugin.CLIMetadata, subcmdMeta *plugin.SubcommandMetadata) {
	subcmdMeta.Description = `Switch the multicluster provider used in cmd/main.go.

Rewrites cmd/main.go while preserving all +kubebuilder:scaffold markers so that
future kubebuilder create api and create webhook commands still work.`
	subcmdMeta.Examples = `  # Switch to namespace provider
  kubebuilder edit --plugins multicluster-runtime.sigs.k8s.io/v1 --provider namespace`
}

func (p *editSubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&p.provider, "provider", "kubeconfig",
		fmt.Sprintf("Switch the multicluster provider (%s)", strings.Join(validProviders, "|")))
	fs.StringVar(&p.kubeconfigDir, "kubeconfig-dir", "/etc/kubeconfig",
		"Directory of per-cluster kubeconfig files (file provider only)")
}

func (p *editSubcommand) InjectConfig(c config.Config) error {
	p.config = c
	return nil
}

func (p *editSubcommand) Scaffold(fs machinery.Filesystem) error {
	if err := validateProvider(p.provider); err != nil {
		return err
	}
	s := scaffolds.NewEditScaffolder(p.config, p.provider, p.kubeconfigDir)
	s.InjectFS(fs)
	return s.Scaffold()
}

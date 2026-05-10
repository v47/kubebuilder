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

var _ plugin.InitSubcommand = &initSubcommand{}

// validProviders is the set of accepted --provider values.
var validProviders = []string{"kubeconfig", "namespace", "cluster-api", "file"}

type initSubcommand struct {
	config        config.Config
	provider      string
	kubeconfigDir string
}

func (p *initSubcommand) UpdateMetadata(_ plugin.CLIMetadata, subcmdMeta *plugin.SubcommandMetadata) {
	subcmdMeta.Description = `Rewrites cmd/main.go to use sigs.k8s.io/multicluster-runtime instead of
the standard single-cluster controller-runtime manager.

Must be chained after go/v4:
  kubebuilder init --plugins go/v4,multicluster-runtime.sigs.k8s.io/v1 ...

The --provider flag selects the cluster-discovery mechanism:
  kubeconfig   Watch kubeconfig Secrets to register clusters at runtime (default)
  namespace    Treat each namespace as a separate "cluster"
  cluster-api  Discover clusters managed by Cluster API
  file         Load clusters from a directory of kubeconfig files`
	subcmdMeta.Examples = `  # Kubeconfig provider (default)
  kubebuilder init --plugins go/v4,multicluster-runtime.sigs.k8s.io/v1 \
    --domain example.com --provider kubeconfig

  # Namespace provider
  kubebuilder init --plugins go/v4,multicluster-runtime.sigs.k8s.io/v1 \
    --domain example.com --provider namespace`
}

func (p *initSubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&p.provider, "provider", "kubeconfig",
		fmt.Sprintf("Multicluster provider (%s)", strings.Join(validProviders, "|")))
	fs.StringVar(&p.kubeconfigDir, "kubeconfig-dir", "/etc/kubeconfig",
		"Directory of per-cluster kubeconfig files (file provider only)")
}

func (p *initSubcommand) InjectConfig(c config.Config) error {
	p.config = c
	return nil
}

func (p *initSubcommand) Scaffold(fs machinery.Filesystem) error {
	if err := validateProvider(p.provider); err != nil {
		return err
	}
	s := scaffolds.NewInitScaffolder(p.config, p.provider, p.kubeconfigDir)
	s.InjectFS(fs)
	return s.Scaffold()
}

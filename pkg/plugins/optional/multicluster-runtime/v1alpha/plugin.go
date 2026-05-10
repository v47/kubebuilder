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

// Package v1 implements the multicluster-runtime/v1 plugin for Kubebuilder.
//
// This plugin modifies the scaffolded project to use sigs.k8s.io/multicluster-runtime
// instead of the standard single-cluster controller-runtime manager, enabling
// controllers to reconcile objects across multiple Kubernetes clusters.
//
// It is designed to be chained after go/v4:
//
//	kubebuilder init --plugins go/v4,multicluster-runtime.sigs.k8s.io/v1 ...
package v1alpha

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	cfgv3 "sigs.k8s.io/kubebuilder/v4/pkg/config/v3"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/stage"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

const (
	// PluginName is the fully qualified plugin name.
	PluginName = "multicluster-runtime.sigs.k8s.io"

	// MulticlusterRuntimeVersion is the version of multicluster-runtime to scaffold against.
	MulticlusterRuntimeVersion = "v0.0.0"
)

var (
	pluginVersion            = plugin.Version{Number: 1, Stage: stage.Alpha}
	supportedProjectVersions = []config.Version{cfgv3.Version}
)

var _ plugin.Init = Plugin{}
var _ plugin.CreateAPI = Plugin{}
var _ plugin.Edit = Plugin{}

// Plugin implements plugin.Init, plugin.CreateAPI, and plugin.Edit.
type Plugin struct {
	initSubcommand
	createAPISubcommand
	editSubcommand
}

// Name returns the plugin's qualified name.
func (Plugin) Name() string { return PluginName }

// Version returns the plugin version.
func (Plugin) Version() plugin.Version { return pluginVersion }

// SupportedProjectVersions returns the project config versions supported by this plugin.
func (Plugin) SupportedProjectVersions() []config.Version { return supportedProjectVersions }

// GetInitSubcommand returns the init subcommand.
func (p Plugin) GetInitSubcommand() plugin.InitSubcommand { return &p.initSubcommand }

// GetCreateAPISubcommand returns the create api subcommand.
func (p Plugin) GetCreateAPISubcommand() plugin.CreateAPISubcommand { return &p.createAPISubcommand }

// GetEditSubcommand returns the edit subcommand.
func (p Plugin) GetEditSubcommand() plugin.EditSubcommand { return &p.editSubcommand }

// DeprecationWarning returns empty — this plugin is not deprecated.
func (Plugin) DeprecationWarning() string { return "" }

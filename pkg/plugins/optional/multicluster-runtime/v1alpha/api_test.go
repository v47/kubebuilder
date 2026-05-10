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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	cfgv3 "sigs.k8s.io/kubebuilder/v4/pkg/config/v3"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
)

var _ = Describe("createAPISubcommand", func() {
	var (
		subCmd *createAPISubcommand
		cfg    config.Config
	)

	BeforeEach(func() {
		subCmd = &createAPISubcommand{}
		cfg = cfgv3.New()
		_ = cfg.SetRepository("github.com/example/myop")
		_ = cfg.SetDomain("example.com")
		Expect(subCmd.InjectConfig(cfg)).To(Succeed())
	})

	Context("InjectConfig", func() {
		It("should store the config", func() {
			Expect(subCmd.config).To(Equal(cfg))
		})
	})

	Context("InjectResource", func() {
		It("should store the resource", func() {
			res := &resource.Resource{
				GVK: resource.GVK{Group: "foo", Version: "v1", Kind: "Bar"},
			}
			Expect(subCmd.InjectResource(res)).To(Succeed())
			Expect(subCmd.resource).To(Equal(res))
		})
	})

	Context("Scaffold", func() {
		var memFS machinery.Filesystem

		BeforeEach(func() {
			memFS = machinery.Filesystem{FS: afero.NewMemMapFs()}
		})

		It("should be a no-op when resource is nil", func() {
			subCmd.resource = nil
			Expect(subCmd.Scaffold(memFS)).To(Succeed())
		})

		It("should be a no-op when resource has no controller", func() {
			subCmd.resource = &resource.Resource{
				GVK:        resource.GVK{Group: "foo", Version: "v1", Kind: "Bar"},
				Controller: false,
			}
			Expect(subCmd.Scaffold(memFS)).To(Succeed())
		})

		It("should write a controller when resource has a controller", func() {
			subCmd.resource = &resource.Resource{
				GVK:        resource.GVK{Group: "foo", Version: "v1", Kind: "Bar"},
				Plural:     "bars",
				Path:       "github.com/example/myop/api/v1",
				Controller: true,
			}
			Expect(subCmd.Scaffold(memFS)).To(Succeed())
			// Verify the controller file was created
			exists, err := afero.Exists(memFS.FS, "internal/controller/bar_controller.go")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})
	})
})

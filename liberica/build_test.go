/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package liberica_test

import (
	"os"
	"testing"

	"github.com/paketo-buildpacks/bellsoft-liberica/liberica"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/libjvm"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it("contributes JDK", func() {
		ctx.Plan.Entries = append(
			ctx.Plan.Entries,
			libcnb.BuildpackPlanEntry{Name: "jdk"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "jdk",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := libjvm.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("jdk"))

		Expect(result.BOM.Entries).To(HaveLen(1))
		Expect(result.BOM.Entries[0].Name).To(Equal("jdk"))
		Expect(result.BOM.Entries[0].Launch).To(BeFalse())
		Expect(result.BOM.Entries[0].Build).To(BeTrue())
	})

	it("contributes NIK", func() {
		ctx.Plan.Entries = append(
			ctx.Plan.Entries,
			libcnb.BuildpackPlanEntry{Name: "jdk", Metadata: map[string]interface{}{}},
			libcnb.BuildpackPlanEntry{Name: "native-image-builder"},
		)
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "native-image-svm",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := liberica.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("native-image-svm"))

		Expect(result.BOM.Entries).To(HaveLen(1))
		Expect(result.BOM.Entries[0].Name).To(Equal("native-image-svm"))
		Expect(result.BOM.Entries[0].Launch).To(BeFalse())
		Expect(result.BOM.Entries[0].Build).To(BeTrue())
	})

	it("contributes NIK alternative buildplan", func() {
		// NIK includes a JDK, so we don't need a second JDK
		ctx.Plan.Entries = append(
			ctx.Plan.Entries,
			libcnb.BuildpackPlanEntry{Name: "native-image-builder"},
			libcnb.BuildpackPlanEntry{Name: "jdk", Metadata: map[string]interface{}{}},
			libcnb.BuildpackPlanEntry{Name: "jre", Metadata: map[string]interface{}{}})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "native-image-svm",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := liberica.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("native-image-svm"))

		Expect(result.BOM.Entries).To(HaveLen(1))
		Expect(result.BOM.Entries[0].Name).To(Equal("native-image-svm"))
		Expect(result.BOM.Entries[0].Launch).To(BeFalse())
		Expect(result.BOM.Entries[0].Build).To(BeTrue())
	})

	it("contributes JRE", func() {
		ctx.Plan.Entries = append(
			ctx.Plan.Entries,
			libcnb.BuildpackPlanEntry{Name: "jre", Metadata: LaunchContribution})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "jre",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := libjvm.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(3))
		Expect(result.Layers[0].Name()).To(Equal("jre"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[2].Name()).To(Equal("java-security-properties"))

		Expect(result.BOM.Entries).To(HaveLen(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("jre"))
		Expect(result.BOM.Entries[0].Launch).To(BeTrue())
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
		Expect(result.BOM.Entries[1].Launch).To(BeTrue())
	})

	it("contributes security-providers-classpath-8 before Java 9", func() {
		ctx.Plan.Entries = append(
			ctx.Plan.Entries,
			libcnb.BuildpackPlanEntry{Name: "jre", Metadata: LaunchContribution})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "jre",
					"version": "8.0.0",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := libjvm.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{
			"active-processor-count",
			"java-opts",
			"jvm-heap",
			"link-local-dns",
			"memory-calculator",
			"openssl-certificate-loader",
			"security-providers-configurer",
			"jmx",
			"security-providers-classpath-8",
			"debug-8",
		}))
	})

	it("contributes security-providers-classpath-9 after Java 9", func() {
		ctx.Plan.Entries = append(
			ctx.Plan.Entries,
			libcnb.BuildpackPlanEntry{Name: "jre", Metadata: LaunchContribution})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "jre",
					"version": "11.0.0",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := libjvm.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{
			"active-processor-count",
			"java-opts",
			"jvm-heap",
			"link-local-dns",
			"memory-calculator",
			"openssl-certificate-loader",
			"security-providers-configurer",
			"jmx",
			"security-providers-classpath-9",
			"debug-9",
			"nmt",
		}))
	})

	it("contributes JDK when no JRE and only a JRE is wanted", func() {
		ctx.Plan.Entries = append(
			ctx.Plan.Entries,
			libcnb.BuildpackPlanEntry{Name: "jre", Metadata: LaunchContribution})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "jdk",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := libjvm.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers[0].Name()).To(Equal("jdk"))
		Expect(result.Layers[0].(libjvm.JRE).LayerContributor.Dependency.ID).To(Equal("jdk"))

		Expect(result.BOM.Entries).To(HaveLen(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("jdk"))
		Expect(result.BOM.Entries[0].Launch).To(BeTrue())
		Expect(result.BOM.Entries[0].Build).To(BeTrue())
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
		Expect(result.BOM.Entries[1].Launch).To(BeTrue())
	})

	it("contributes JDK when no JRE and both a JDK and JRE are wanted", func() {
		ctx.Plan.Entries = append(
			ctx.Plan.Entries,
			libcnb.BuildpackPlanEntry{Name: "jdk", Metadata: LaunchContribution},
			libcnb.BuildpackPlanEntry{Name: "jre", Metadata: LaunchContribution})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "jdk",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := libjvm.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers[0].Name()).To(Equal("jdk"))
		Expect(result.Layers[0].(libjvm.JRE).LayerContributor.Dependency.ID).To(Equal("jdk"))

		Expect(result.BOM.Entries).To(HaveLen(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("jdk"))
		Expect(result.BOM.Entries[0].Launch).To(BeTrue())
		Expect(result.BOM.Entries[0].Build).To(BeTrue())
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
		Expect(result.BOM.Entries[1].Launch).To(BeTrue())
	})

	context("$BP_JVM_VERSION", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_JVM_VERSION", "1.1.1")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_JVM_VERSION")).To(Succeed())
		})

		it("selects versions based on BP_JVM_VERSION", func() {
			ctx.Plan.Entries = append(ctx.Plan.Entries,
				libcnb.BuildpackPlanEntry{Name: "jdk"},
				libcnb.BuildpackPlanEntry{Name: "jre"},
			)
			ctx.Buildpack.Metadata = map[string]interface{}{
				"dependencies": []map[string]interface{}{
					{
						"id":      "jdk",
						"version": "1.1.1",
						"stacks":  []interface{}{"test-stack-id"},
					},
					{
						"id":      "jdk",
						"version": "2.2.2",
						"stacks":  []interface{}{"test-stack-id"},
					},
					{
						"id":      "jre",
						"version": "1.1.1",
						"stacks":  []interface{}{"test-stack-id"},
					},
					{
						"id":      "jre",
						"version": "2.2.2",
						"stacks":  []interface{}{"test-stack-id"},
					},
				},
			}
			ctx.StackID = "test-stack-id"

			result, err := libjvm.Build{}.Build(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Layers[0].(libjvm.JDK).LayerContributor.Dependency.Version).To(Equal("1.1.1"))
			Expect(result.Layers[1].(libjvm.JRE).LayerContributor.Dependency.Version).To(Equal("1.1.1"))
		})
	})

	context("$BP_JVM_TYPE", func() {

		it.After(func() {
			Expect(os.Unsetenv("BP_JVM_TYPE")).To(Succeed())
		})

		it("contributes JDK when specified explicitly in $BP_JVM_TYPE", func() {
			Expect(os.Setenv("BP_JVM_TYPE", "jdk")).To(Succeed())

			ctx.Plan.Entries = append(
				ctx.Plan.Entries,
				libcnb.BuildpackPlanEntry{Name: "jdk", Metadata: LaunchContribution},
				libcnb.BuildpackPlanEntry{Name: "jre", Metadata: LaunchContribution})
			ctx.Buildpack.Metadata = map[string]interface{}{
				"dependencies": []map[string]interface{}{
					{
						"id":      "jdk",
						"version": "0.0.2",
						"stacks":  []interface{}{"test-stack-id"},
					},
					{
						"id":      "jre",
						"version": "2.2.2",
						"stacks":  []interface{}{"test-stack-id"},
					},
				},
			}
			ctx.StackID = "test-stack-id"

			result, err := libjvm.Build{}.Build(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Layers[0].Name()).To(Equal("jdk"))
			Expect(result.Layers[0].(libjvm.JRE).LayerContributor.Dependency.ID).To(Equal("jdk"))

			Expect(result.BOM.Entries[0].Name).To(Equal("jdk"))
			Expect(result.BOM.Entries[0].Launch).To(BeTrue())
		})

		it("contributes JRE when specified explicitly in $BP_JVM_TYPE", func() {
			Expect(os.Setenv("BP_JVM_TYPE", "jre")).To(Succeed())

			ctx.Plan.Entries = append(
				ctx.Plan.Entries,
				libcnb.BuildpackPlanEntry{Name: "jdk", Metadata: LaunchContribution},
				libcnb.BuildpackPlanEntry{Name: "jre", Metadata: LaunchContribution})
			ctx.Buildpack.Metadata = map[string]interface{}{
				"dependencies": []map[string]interface{}{
					{
						"id":      "jdk",
						"version": "0.0.1",
						"stacks":  []interface{}{"test-stack-id"},
					},
					{
						"id":      "jre",
						"version": "1.1.1",
						"stacks":  []interface{}{"test-stack-id"},
					},
				},
			}
			ctx.StackID = "test-stack-id"

			result, err := libjvm.Build{}.Build(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Layers[0].Name()).To(Equal("jdk"))
			Expect(result.Layers[0].(libjvm.JDK).LayerContributor.Dependency.ID).To(Equal("jdk"))
			Expect(result.Layers[1].(libjvm.JRE).LayerContributor.Dependency.ID).To(Equal("jre"))

			Expect(result.BOM.Entries).To(HaveLen(3))
			Expect(result.BOM.Entries[0].Name).To(Equal("jdk"))
			Expect(result.BOM.Entries[0].Launch).To(BeFalse())
			Expect(result.BOM.Entries[0].Build).To(BeTrue())

			Expect(result.BOM.Entries[1].Name).To(Equal("jre"))
			Expect(result.BOM.Entries[1].Launch).To(BeTrue())
			Expect(result.BOM.Entries[1].Build).To(BeTrue())

		})
	})
}

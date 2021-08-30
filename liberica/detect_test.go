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
	"testing"

	"github.com/paketo-buildpacks/bellsoft-liberica/liberica"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect liberica.Detect
	)

	it("includes build plan options", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jdk"},
						{Name: "native-image-builder"},
						{Name: "jre"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "jdk"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jdk"},
						{Name: "native-image-builder"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "jdk"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jdk"},
						{Name: "jre"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jdk"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jre"},
					},
				},
			},
		}))

	})
}

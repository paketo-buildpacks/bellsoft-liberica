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

package main

import (
	"github.com/paketo-buildpacks/libjvm"
	"github.com/paketo-buildpacks/libjvm/helper"
	"os"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

func main() {
	logger := bard.NewLogger(os.Stdout)

	// not used directly, but this forces the helper module to be included in the module
	// we need the helper module because of the way that `scripts/build/sh` builds the helper cmd
	_ = helper.ActiveProcessorCount{Logger: logger}

	libpak.Main(
		libjvm.Detect{},
		libjvm.NewBuild(logger, libjvm.WithNativeImage(
			libjvm.NativeImage{
				BundledWithJDK: true,
			}),
			libjvm.WithCustomHelpers([]string{"nmt"})),
	)
}

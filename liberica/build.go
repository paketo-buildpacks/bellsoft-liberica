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

package liberica

import (
	"fmt"
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libjvm"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"io/ioutil"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	build := libjvm.NewBuild(b.Logger)

	var err error
	var version string
	var nativeImage bool

	if _, nativeImage, err = pr.Resolve("native-image-builder"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve native-image-builder plan entry\n%w", err)
	}
	version, err = b.getJavaVersion(context, nativeImage)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to get Java version\n%w", err)
	}

	if nativeImage {
		dr, err := libpak.NewDependencyResolver(context)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
		}
		dep, err := dr.Resolve("native-image-svm", version)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}
		jdk, be, err := libjvm.NewJDK(dep, build.DependencyCache, build.CertLoader)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create jdk\n%w", err)
		}

		jdk.Logger = build.Logger
		build.Result.Layers = append(build.Result.Layers, jdk)
		build.Result.BOM.Entries = append(build.Result.BOM.Entries, be)
		return build.Result, nil
	}

	r, err := build.Build(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to build\n%w", err)
	}
	if libjvm.IsBeforeJava9(version) {
		h, be := libpak.NewHelperLayer(context.Buildpack, "nmt")
		h.Logger = build.Logger
		r.Layers = append(r.Layers, h)
		r.BOM.Entries = append(r.BOM.Entries, be)
	}
	return r, nil
}

func (b *Build) getJavaVersion(context libcnb.BuildContext, shouldLog bool) (string, error) {
	var logger bard.Logger
	if shouldLog {
		b.Logger.Title(context.Buildpack)
		logger = b.Logger
	} else {
		logger = bard.NewLogger(ioutil.Discard)
	}

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &logger)
	if err != nil {
		return "", fmt.Errorf("unable to create configuration resolver\n%w", err)
	}
	jvmVersion := libjvm.NewJVMVersion(logger)
	v, err := jvmVersion.GetJVMVersion(context.Application.Path, cr)
	if err != nil {
		return "", fmt.Errorf("unable to determine jvm version\n%w", err)
	}
	return v, nil
}

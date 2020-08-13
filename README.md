# `gcr.io/paketo-buildpacks/bellsoft-liberica`
The Paketo BellSoft Liberica Buildpack is a Cloud Native Buildpack that provides the BellSoft Liberica implementations of JREs and JDKs.

This buildpack is designed to work in collaboration with other buildpacks which request contributions of JREs and JDKs.

## Behavior
This buildpack will participate if any of the following conditions are met

* Another buildpack requires `jdk`
* Another buildpack requires `jre`

The buildpack will do the following if a JDK is requested:

* Contributes a JDK to a layer marked `build` and `cache` with all commands on `$PATH`
* Contributes `$JAVA_HOME` configured to the build layer
* Contributes `$JDK_HOME` configure to the build layer

The buildpack will do the following if a JRE is requested:

* Contributes a JRE to a layer with all commands on `$PATH`
* Contributes `$JAVA_HOME` configured to the layer
* Contributes `-XX:ActiveProcessorCount` to the layer
* Contributes `$MALLOC_ARENA_MAX` to the layer
* Disables JVM DNS caching if link-local DNS is available
* If `metadata.build = true`
  * Marks layer as `build` and `cache`
* If `metadata.launch = true`
  * Marks layer as `launch`
* Contributes `jvmkill` to a layer marked `launch`
* Contributes Memory Calculator to a layer marked `launch`

## Configuration
| Environment Variable | Description
| -------------------- | -----------
| `$BP_JVM_VERSION` | Configure a specific JDK or JRE version.  This value must _exactly_ match a version available in the buildpack so, it's best to use a wildcard such as `8.*`.<br/>  If you would like to to pin the exact version used, consider pinning your builder image, e.g. instead of `gcr.io/paketo-buildpacks/builder:base` use `gcr.io/paketo-buildpacks/builder:0.0.389-base`. The latest version can be found via the web interface of GCR, e.g. [here](https://gcr.io/paketo-buildpacks/builder). The corresponding bellsoft-liberica version is listed the builder image's label `io.buildpacks.builder.metadata`, so either pull the image and `docker image inspect` or use [`skopeo inspect`](https://github.com/containers/skopeo). The JDK & JRE versions corresponding to the buildpack versions can be found in [releases](https://github.com/paketo-buildpacks/bellsoft-liberica/releases).
| `$BPL_JVM_HEAD_ROOM` | Configure the percentage of headroom the memory calculator will allocated.  Defaults to `0`.
| `$BPL_JVM_LOADED_CLASS_COUNT` | Configure the number of classes that will be loaded at runtime.  Defaults to 35% of the number of classes.
| `$BPL_JVM_THREAD_COUNT` | Configure the number of user threads at runtime.  Defaults to `250`.
| `$JAVA_TOOL_OPTIONS` | Configure the JVM launch flags

## Bindings
The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`
|Key                   | Value   | Description
|----------------------|---------|------------
|`<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>`

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

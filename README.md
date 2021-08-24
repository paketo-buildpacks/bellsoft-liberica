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
| `$BP_JVM_VERSION` | Configure the JVM version (e.g. `8`, `11`, `16`).  The buildpack will download JDK and JRE assets that are compatible with this version of the JVM specification.  Since the buildpack only ships a single version of each supported line, updates to the buildpack can change the exact version of the JDK or JRE.  In order to hold the JDK and JRE versions stable, the buildpack version itself must be stable.<p/><p/>Buildpack releases (and the dependency versions for each release) can be found [here][bpv].  Few users will use this buildpack directly, instead consuming a language buildpack like `paketo-buildpacks/java` who's releases (and the individual buildpack versions and dependency versions for each release) can be found [here](https://github.com/paketo-buildpacks/java/releases).  Finally, some users will will consume builders like `paketobuildpacks/builder:base` who's releases can be found [here](https://hub.docker.com/r/paketobuildpacks/builder/tags?page=1&name=base).  To determine the individual buildpack versions and dependency versions for each builder release use the [`pack inspect-builder <image>`](https://buildpacks.io/docs/reference/pack/pack_inspect-builder/) functionality.
| `$BP_JVM_TYPE` | Configure the JVM type that is provided at runtime, i.e. a JDK or JRE - accepts values "JDK" or "JRE" (default). If a JRE type is requested but not available, a JDK will be provided.
| `$BPL_JVM_HEAD_ROOM` | Configure the percentage of headroom the memory calculator will allocated.  Defaults to `0`.
| `$BPL_JVM_LOADED_CLASS_COUNT` | Configure the number of classes that will be loaded at runtime.  Defaults to 35% of the number of classes.
| `$BPL_JVM_THREAD_COUNT` | Configure the number of user threads at runtime.  Defaults to `250`.
| `$JAVA_TOOL_OPTIONS` | Configure the JVM launch flags

[bpv]: https://github.com/paketo-buildpacks/bellsoft-liberica/releases

## Bindings
The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`
|Key                   | Value   | Description
|----------------------|---------|------------
|`<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>`

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0


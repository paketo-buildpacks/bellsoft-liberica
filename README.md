# `gcr.io/paketo-buildpacks/bellsoft-liberica`

The Paketo BellSoft Liberica Buildpack is a Cloud Native Buildpack that provides the BellSoft Liberica implementations of JREs and JDKs.

This buildpack is designed to work in collaboration with other buildpacks which request contributions of JREs and JDKs.

## Integration

Downstream buildpacks can require the JRE/JDK dependency by generating a [Build Plan
TOML](https://github.com/buildpacks/spec/blob/master/buildpack.md#build-plan-toml)
file that looks like the following:

```toml
[[requires]]

  # The name of the JRE dependency is "jre". This value is considered
  # part of the public API for the buildpack and will not change without a plan
  # for deprecation.
  name = "jre"

  # The version of the JRE dependency is not required. In the case it
  # is not specified, the buildpack will provide the default version, which can
  # be seen in the buildpack.toml file.
  # If you wish to request a specific version, the buildpack supports
  # specifying a semver constraint in the form of "11.*", "11.0.*", or even
  # "11.0.13".
  version = "11.0.13"

  # The BellSoft Liberica buildpack supports some non-required metadata options.
  [requires.metadata]
    # Setting the launch flag to true will ensure that the BellSoft Liberica
    # dependency is available on the $PATH for the running application. If you are
    # writing an application that needs to run node at runtime, this flag should
    # be set to true.
    launch = true

[[requires]]

  # The name of the JDK dependency is "jdk". This value is considered
  # part of the public API for the buildpack and will not change without a plan
  # for deprecation.
  name = "jdk"

  # The version of the JDK dependency is not required. In the case it
  # is not specified, the buildpack will provide the default version, which can
  # be seen in the buildpack.toml file.
  # If you wish to request a specific version, the buildpack supports
  # specifying a semver constraint in the form of "11.*", "11.0.*", or even
  # "11.0.13".
  version = "11.0.13"
```

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
* Contributes `-XX:+ExitOnOutOfMemoryError` to the layer
* Contributes `-XX:+UnlockDiagnosticVMOptions`,`-XX:NativeMemoryTracking=summary` & `-XX:+PrintNMTStatistics` to the layer (Java NMT)
* If `BPL_JMX_ENABLED = true`
  * Contributes `-Djava.rmi.server.hostname=127.0.0.1`, `-Dcom.sun.management.jmxremote.authenticate=false`, `-Dcom.sun.management.jmxremote.ssl=false` & `-Dcom.sun.management.jmxremote.rmi.port=5000`
* If `BPL_DEBUG_ENABLED = true`
  * Contributes `-agentlib:jdwp=transport=dt_socket,server=y,address=*:8000,suspend=n`. If Java version is 8, address parameter is `address=:8000`
* If `BPL_JFR_ENABLED = true`
  * Contributes `-XX:StartFlightRecording=dumponexit=true,filename=/tmp/recording.jfr`
* Contributes `$MALLOC_ARENA_MAX` to the layer
* Disables JVM DNS caching if link-local DNS is available
* If `metadata.build = true`
  * Marks layer as `build` and `cache`
* If `metadata.launch = true`
  * Marks layer as `launch`
* Contributes Memory Calculator to a layer marked `launch`
* Contributes Heap Dump helper to a layer marked `launch`

## Configuration

| Environment Variable          | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| ----------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `$BP_JVM_VERSION`             | Configure the JVM version (e.g. `8`, `11`, `16`).  The buildpack will download JDK and JRE assets that are compatible with this version of the JVM specification.  Since the buildpack only ships a single version of each supported line, updates to the buildpack can change the exact version of the JDK or JRE.  In order to hold the JDK and JRE versions stable, the buildpack version itself must be stable.<p/><p/>Buildpack releases (and the dependency versions for each release) can be found [here][bpv].  Few users will use this buildpack directly, instead consuming a language buildpack like `paketo-buildpacks/java` who's releases (and the individual buildpack versions and dependency versions for each release) can be found [here](https://github.com/paketo-buildpacks/java/releases).  Finally, some users will will consume builders like `paketobuildpacks/builder:base` who's releases can be found [here](https://hub.docker.com/r/paketobuildpacks/builder/tags?page=1&name=base).  To determine the individual buildpack versions and dependency versions for each builder release use the [`pack inspect-builder <image>`](https://buildpacks.io/docs/reference/pack/pack_inspect-builder/) functionality. |
| `$BP_JVM_TYPE`                | Configure the JVM type that is provided at runtime, i.e. a JDK or JRE - accepts values "JDK" or "JRE" (default). If a JRE type is requested but not available, a JDK will be provided.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| `$BPL_JVM_HEAD_ROOM`          | Configure the percentage of headroom the memory calculator will allocated.  Defaults to `0`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `$BPL_JVM_LOADED_CLASS_COUNT` | Configure the number of classes that will be loaded at runtime.  Defaults to 35% of the number of classes.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| `$BPL_JVM_THREAD_COUNT`       | Configure the number of user threads at runtime.  Defaults to `250`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `$BPL_HEAP_DUMP_PATH`         | Configure the location for writing heap dumps in the event of an OutOfMemoryError exception. Defaults to ``, which disables writing heap dumps. The path set must be writable by the JVM process.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `$BPL_JAVA_NMT_ENABLED`         | Configure whether Java Native Memory Tracking (NMT) is enabled. Defaults to `true`. Set this to `false` to disable NMT functionality.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `$BPL_JAVA_NMT_LEVEL`         | Configure the level of detail for Java Native Memory Tracking (NMT) output. Defaults to `summary`. Set this to `detail` for detailed NMT output.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `$BPL_JMX_ENABLED`         | Configure whether Java Management Extensions (JMX) is enabled. Defaults to `false`. Set this to `true` to enable JMX functionality.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `$BPL_JMX_PORT`         | Configure the port number for JMX. Defaults to `5000`. When running the container, this value should match the port published locally, i.e. for Docker: --publish 5000:5000                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| `$BPL_DEBUG_ENABLED`         | Configure whether remote debugging features are enabled. Defaults to `false`. Set this to `true` to enable remote debugging.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `$BPL_DEBUG_PORT`         | Configure the port number for remote debugging. Defaults to `8000`.
| `$BPL_DEBUG_SUSPEND`         | Configure whether to suspend execution until a debugger has attached. Defaults to `false`.
| `$BPL_JFR_ENABLED`         | Configure whether Java Flight Recording (JFR) is enabled. If no arguments are specified via `BPL_JFR_ARGS`, the default config args `dumponexit=true,filename=/tmp/recording.jfr` are added.
| `$BPL_JFR_ARGS`         | Configure custom arguments to Java Flight Recording, via a comma-separated list, e.g. `duration=10s,maxage=1m`. If any values are specified, no default args are supplied.
| `$JAVA_TOOL_OPTIONS`          | Configure the JVM launch flags                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |

[bpv]: https://github.com/paketo-buildpacks/bellsoft-liberica/releases

## Bindings

The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`

| Key                   | Value   | Description                                                                                       |
| --------------------- | ------- | ------------------------------------------------------------------------------------------------- |
| `<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>` |

## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0


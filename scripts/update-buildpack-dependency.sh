sha256() {
  if [[ "${DEPENDENCY}" == "jvmkill" ]]; then
    shasum -a 256 "${ROOT}"/dependency/jvmkill-*.so | cut -f 1 -d ' '
  elif [[ "${DEPENDENCY}" == "memory-calculator" ]]; then
    shasum -a 256 "${ROOT}"/dependency/memory-calculator-*.tgz | cut -f 1 -d ' '
  else
    cat "${ROOT}"/dependency/sha256
  fi
}

uri() {
  if [[ "${DEPENDENCY}" == "jvmkill" ]]; then
    echo "https://github.com/cloudfoundry/jvmkill/releases/download/v$(cat "${ROOT}"/dependency/version)/$(basename "${ROOT}"/dependency/jvmkill-*.so)"
  elif [[ "${DEPENDENCY}" == "memory-calculator" ]]; then
    echo "https://github.com/cloudfoundry/java-buildpack-memory-calculator/releases/download/v$(cat "${ROOT}"/dependency/version)/$(basename "${ROOT}"/dependency/memory-calculator-*.tgz)"
  else
    cat "${ROOT}"/dependency/uri
  fi
}

version() {
  local PATTERN="([0-9]+)\.([0-9]+)\.([0-9]+)(.*)"

  if [[ $(cat "${ROOT}"/dependency/version) =~ ${PATTERN} ]]; then
    echo "${BASH_REMATCH[1]}.${BASH_REMATCH[2]}.${BASH_REMATCH[3]}"
    return
  else
    echo "version is not semver" 1>&2
    exit 1
  fi
}

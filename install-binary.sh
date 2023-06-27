#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

trap_add() {
  local -r sig="${2:?Signal required}"
  local -r hdls="$(trap -p "${sig}" | cut -f2 -d \')"
  # shellcheck disable=SC2064 # Quotes are required here to properly expand when adding the new trap.
  trap "${hdls}${hdls:+;}${1:?Handler required}" "${sig}"
}

OS_NAME="$(uname -s)"
readonly OS_NAME
OS_ARCH="$(uname -m)"
readonly OS_ARCH

declare -r PLUGIN_NAME="helm-list-images"
VERSION="$(grep VERSION "${HELM_PLUGIN_DIR}/plugin.yaml" | cut -d'"' -f2)"
readonly VERSION
declare -r DOWNLOAD_URL="https://github.com/d2iq-labs/${PLUGIN_NAME}/releases/download/v${VERSION}/${PLUGIN_NAME}_${VERSION}_${OS_NAME}_${OS_ARCH}.tar.gz"

echo -e "download url set to ${DOWNLOAD_URL}\n"
echo -e "artifact name with path ${OUTPUT_BASENAME_WITH_POSTFIX}\n"
echo -e "downloading ${DOWNLOAD_URL} to ${HELM_PLUGIN_DIR}\n"

if [ -z "${DOWNLOAD_URL}" ]; then
  echo -e "Unsupported OS / architecture: ${OS_NAME}/${OS_ARCH}\n"
  exit 1
fi

HELM_PLUGIN_TEMP_PATH="$(mktemp -d -p "${TMPDIR:-/tmp}" "${PLUGIN_NAME}_XXXXXXX")"
readonly HELM_PLUGIN_TEMP_PATH
trap_add "rm -rf \"${HELM_PLUGIN_TEMP_PATH}\"" EXIT

if command -v curl &>/dev/null; then
  if curl -fsSL "${DOWNLOAD_URL}" | tar xz -C "${HELM_PLUGIN_TEMP_PATH}"; then
    echo -e "successfully downloaded and extracted the plugin archive\n"
  else
    echo -e "failed while downloading helm archive\n"
    exit 1
  fi
else
  echo "Need curl"
  exit 1
fi

mkdir -p "${HELM_PLUGIN_DIR}/bin"
mv "${HELM_PLUGIN_TEMP_PATH}/${PLUGIN_NAME}" "${HELM_PLUGIN_DIR}/bin/${PLUGIN_NAME}"

echo
echo "${PLUGIN_NAME} is installed."
echo
"${HELM_PLUGIN_DIR}/bin/${PLUGIN_NAME}" -h
echo
echo "See https://github.com/d2iq-labs/${PLUGIN_NAME}#readme for more information on getting started."

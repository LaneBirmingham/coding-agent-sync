#!/usr/bin/env bash
set -euo pipefail

REPO="LaneBirmingham/coding-agent-sync"
VERSION="${CAS_VERSION:-latest}"
INSTALL_DIR="${CAS_INSTALL_DIR:-${HOME}/.local/bin}"
BIN_NAME="${CAS_BIN_NAME:-cas}"

usage() {
  cat <<'USAGE'
Install coding-agent-sync (cas) from GitHub releases.

Usage:
  install.sh [--version <version>] [--install-dir <dir>] [--bin-name <name>]

Environment overrides:
  CAS_VERSION      Version to install, for example 0.3.0 or v0.3.0 (default: latest)
  CAS_INSTALL_DIR  Install directory (default: ~/.local/bin)
  CAS_BIN_NAME     Installed binary name (default: cas)
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      if [[ $# -lt 2 ]]; then
        echo "Missing value for --version" >&2
        exit 1
      fi
      VERSION="${2:-}"
      shift 2
      ;;
    --install-dir)
      if [[ $# -lt 2 ]]; then
        echo "Missing value for --install-dir" >&2
        exit 1
      fi
      INSTALL_DIR="${2:-}"
      shift 2
      ;;
    --bin-name)
      if [[ $# -lt 2 ]]; then
        echo "Missing value for --bin-name" >&2
        exit 1
      fi
      BIN_NAME="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

need_cmd curl
need_cmd tar
need_cmd uname

if [[ -z "${VERSION}" ]]; then
  echo "Version cannot be empty" >&2
  exit 1
fi
if [[ -z "${INSTALL_DIR}" ]]; then
  echo "Install directory cannot be empty" >&2
  exit 1
fi
if [[ -z "${BIN_NAME}" ]]; then
  echo "Binary name cannot be empty" >&2
  exit 1
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "${OS}" in
  darwin|linux) ;;
  *)
    echo "Unsupported OS: ${OS}" >&2
    exit 1
    ;;
esac

case "${ARCH}" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: ${ARCH}" >&2
    exit 1
    ;;
esac

if [[ "${VERSION}" == "latest" ]]; then
  TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": "v\([^"]*\)".*/\1/p')"
  if [[ -z "${TAG}" ]]; then
    echo "Failed to resolve latest release version from GitHub API" >&2
    exit 1
  fi
  VERSION="${TAG}"
else
  VERSION="${VERSION#v}"
fi

ASSET="cas_${VERSION}_${OS}_${ARCH}.tar.gz"
BIN_IN_ARCHIVE="cas_${VERSION}_${OS}_${ARCH}"
BASE_URL="https://github.com/${REPO}/releases/download/v${VERSION}"

tmpdir="$(mktemp -d)"
cleanup() {
  rm -rf "${tmpdir}"
}
trap cleanup EXIT

echo "Downloading ${ASSET}..."
curl -fsSL "${BASE_URL}/${ASSET}" -o "${tmpdir}/${ASSET}"
curl -fsSL "${BASE_URL}/SHA256SUMS" -o "${tmpdir}/SHA256SUMS"

echo "Verifying checksum..."
expected="$(awk -v asset="${ASSET}" '$2 ~ ("(^|/)" asset "$") {print $1; exit}' "${tmpdir}/SHA256SUMS")"
if [[ -z "${expected}" ]]; then
  echo "Could not find checksum entry for ${ASSET} in SHA256SUMS" >&2
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  actual="$(sha256sum "${tmpdir}/${ASSET}" | awk '{print $1}')"
elif command -v shasum >/dev/null 2>&1; then
  actual="$(shasum -a 256 "${tmpdir}/${ASSET}" | awk '{print $1}')"
else
  echo "Missing required checksum tool: sha256sum or shasum" >&2
  exit 1
fi

if [[ "${actual}" != "${expected}" ]]; then
  echo "Checksum verification failed for ${ASSET}" >&2
  exit 1
fi

tar -xzf "${tmpdir}/${ASSET}" -C "${tmpdir}"
mkdir -p "${INSTALL_DIR}"

if command -v install >/dev/null 2>&1; then
  install -m 0755 "${tmpdir}/${BIN_IN_ARCHIVE}" "${INSTALL_DIR}/${BIN_NAME}"
else
  cp "${tmpdir}/${BIN_IN_ARCHIVE}" "${INSTALL_DIR}/${BIN_NAME}"
  chmod 0755 "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "Installed ${BIN_NAME} to ${INSTALL_DIR}/${BIN_NAME}"
"${INSTALL_DIR}/${BIN_NAME}" version

case ":${PATH}:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    echo "Add ${INSTALL_DIR} to PATH if needed:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    ;;
esac

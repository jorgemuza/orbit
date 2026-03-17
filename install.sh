#!/bin/sh
set -e

# orbit CLI installer
# Usage: curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh

REPO="jorgemuza/orbit"
BINARY="orbit"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  mingw*|msys*|cygwin*) OS="windows" ;;
  *) echo "Error: unsupported OS: $OS" >&2; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  aarch64|arm64)  ARCH="arm64" ;;
  *) echo "Error: unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

# Get latest version
VERSION=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
if [ -z "$VERSION" ]; then
  echo "Error: could not determine latest version" >&2
  exit 1
fi

echo "Installing ${BINARY} ${VERSION} (${OS}/${ARCH})..."

# Build download URL
if [ "$OS" = "windows" ]; then
  EXT="zip"
else
  EXT="tar.gz"
fi
FILENAME="${BINARY}_${VERSION#v}_${OS}_${ARCH}.${EXT}"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

# Download and install
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${URL}..."
curl -sSfL "$URL" -o "${TMPDIR}/${FILENAME}"

if [ "$EXT" = "tar.gz" ]; then
  tar -xzf "${TMPDIR}/${FILENAME}" -C "$TMPDIR"
else
  unzip -q "${TMPDIR}/${FILENAME}" -d "$TMPDIR"
fi

# Install binary
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

chmod +x "${INSTALL_DIR}/${BINARY}"

echo "Successfully installed ${BINARY} ${VERSION} to ${INSTALL_DIR}/${BINARY}"

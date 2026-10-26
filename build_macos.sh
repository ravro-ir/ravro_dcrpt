#!/bin/bash

# Build one macOS app that runs on Intel (x86_64) and Apple Silicon (arm64).
#
# OpenSSL is compiled from source for both architectures and linked statically,
# then the two program slices are merged with lipo. The resulting app does not
# need a Homebrew OpenSSL on the Mac that runs it.

set -euo pipefail

if [[ "$(uname -s)" != "Darwin" ]]; then
    echo "This script must be run on macOS."
    exit 1
fi

OPENSSL_VERSION="3.5.1"
OPENSSL_URL="https://github.com/openssl/openssl/releases/download/openssl-${OPENSSL_VERSION}/openssl-${OPENSSL_VERSION}.tar.gz"
DEPLOYMENT_TARGET="13.0"
ROOT="$(cd "$(dirname "$0")" && pwd)"
BUILD_DIR="$ROOT/build"
SRC_DIR="$BUILD_DIR/src"

native_goarch() {
    case "$(uname -m)" in
        arm64) echo arm64 ;;
        x86_64) echo amd64 ;;
        *)
            echo "Unsupported architecture: $(uname -m)" >&2
            exit 1
            ;;
    esac
}

clang_arch() {
    case "$1" in
        arm64) echo arm64 ;;
        amd64) echo x86_64 ;;
        *)
            echo "Unsupported Go arch: $1" >&2
            exit 1
            ;;
    esac
}

openssl_target() {
    case "$1" in
        arm64) echo darwin64-arm64 ;;
        amd64) echo darwin64-x86_64 ;;
        *)
            echo "Unsupported Go arch: $1" >&2
            exit 1
            ;;
    esac
}

fetch_openssl_source() {
    local origin="$SRC_DIR/openssl-${OPENSSL_VERSION}"
    local tarball="$SRC_DIR/openssl-${OPENSSL_VERSION}.tar.gz"

    mkdir -p "$SRC_DIR"
    if [[ ! -f "$tarball" ]]; then
        echo "Downloading OpenSSL ${OPENSSL_VERSION}..."
        curl -fL --retry 3 -o "$tarball" "$OPENSSL_URL"
    fi
    if [[ ! -d "$origin" ]]; then
        tar -xzf "$tarball" -C "$SRC_DIR"
    fi
}

build_static_openssl() {
    local goarch="$1"
    local prefix="$BUILD_DIR/openssl-${goarch}"
    local crypto="$prefix/lib/libcrypto.a"
    local ssl="$prefix/lib/libssl.a"

    if [[ -f "$crypto" && -f "$ssl" ]]; then
        echo "Using existing static OpenSSL for ${goarch}: ${prefix}"
        return 0
    fi

    fetch_openssl_source

    local origin="$SRC_DIR/openssl-${OPENSSL_VERSION}"
    local work="$SRC_DIR/openssl-${OPENSSL_VERSION}-${goarch}"
    local target extra=()
    target="$(openssl_target "$goarch")"
    if [[ "$goarch" != "$(native_goarch)" ]]; then
        extra+=(no-asm)
    fi

    echo "Building static OpenSSL ${OPENSSL_VERSION} for ${goarch} (${target})..."
    rm -rf "$work" "$prefix"
    cp -a "$origin" "$work"
    (
        cd "$work"
        ./Configure "$target" no-shared no-zlib no-tests "${extra[@]+"${extra[@]}"}" \
            --prefix="$prefix" \
            --libdir=lib \
            -mmacosx-version-min="$DEPLOYMENT_TARGET"
        make -j"$(sysctl -n hw.ncpu)" build_libs
    )

    mkdir -p "$prefix/lib" "$prefix/include"
    cp "$work/libssl.a" "$work/libcrypto.a" "$prefix/lib/"
    rm -rf "$prefix/include/openssl"
    cp -R "$work/include/openssl" "$prefix/include/openssl"

    local info
    info="$(lipo -info "$crypto")"
    case "$goarch" in
        arm64) [[ "$info" == *arm64* ]] || { echo "OpenSSL archive is not arm64: $info"; exit 1; } ;;
        amd64) [[ "$info" == *x86_64* ]] || { echo "OpenSSL archive is not x86_64: $info"; exit 1; } ;;
    esac
}

build_slice() {
    local goarch="$1"
    local prefix="$BUILD_DIR/openssl-${goarch}"
    local arch sdk
    arch="$(clang_arch "$goarch")"
    sdk="$(xcrun --sdk macosx --show-sdk-path)"

    echo "Building ravro_dcrpt darwin/${goarch}..."
    env \
        CGO_ENABLED=1 \
        GOOS=darwin \
        GOARCH="$goarch" \
        CC="clang -arch ${arch}" \
        SDKROOT="$sdk" \
        MACOSX_DEPLOYMENT_TARGET="$DEPLOYMENT_TARGET" \
        CGO_CFLAGS="-arch ${arch} -isysroot ${sdk} -mmacosx-version-min=${DEPLOYMENT_TARGET} -I${prefix}/include" \
        CGO_LDFLAGS="-arch ${arch} -isysroot ${sdk} -mmacosx-version-min=${DEPLOYMENT_TARGET} ${prefix}/lib/libssl.a ${prefix}/lib/libcrypto.a -framework Security -framework CoreFoundation" \
        go build -ldflags="-s -w" -o "$BUILD_DIR/ravro_dcrpt-${goarch}" ./cmd/cli

    env \
        CGO_ENABLED=1 \
        GOOS=darwin \
        GOARCH="$goarch" \
        CC="clang -arch ${arch}" \
        SDKROOT="$sdk" \
        MACOSX_DEPLOYMENT_TARGET="$DEPLOYMENT_TARGET" \
        CGO_CFLAGS="-arch ${arch} -isysroot ${sdk} -mmacosx-version-min=${DEPLOYMENT_TARGET} -I${prefix}/include" \
        CGO_LDFLAGS="-arch ${arch} -isysroot ${sdk} -mmacosx-version-min=${DEPLOYMENT_TARGET} ${prefix}/lib/libssl.a ${prefix}/lib/libcrypto.a -framework Security -framework CoreFoundation" \
        go build -ldflags="-s -w" -o "$BUILD_DIR/ravro_dcrpt_gui-${goarch}" ./cmd/gui
}

assert_static_openssl() {
    local binary="$1"
    local arch="$2"
    local libs
    libs="$(otool -arch "$arch" -L "$binary")"
    if grep -E 'libssl|libcrypto|openssl' <<<"$libs" >/dev/null; then
        echo "OpenSSL is still dynamically linked in ${binary} (${arch}):"
        echo "$libs"
        exit 1
    fi
}

write_app_bundle() {
    local binary="$1"
    local app="$BUILD_DIR/Ravro Decryption Tool.app"

    rm -rf "$app"
    mkdir -p "$app/Contents/MacOS" "$app/Contents/Resources"
    cp "$binary" "$app/Contents/MacOS/ravro_dcrpt_gui"
    chmod +x "$app/Contents/MacOS/ravro_dcrpt_gui"

    cat > "$app/Contents/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>ravro_dcrpt_gui</string>
    <key>CFBundleIdentifier</key>
    <string>com.ravro.decryption</string>
    <key>CFBundleName</key>
    <string>Ravro Decryption Tool</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>2.0.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>${DEPLOYMENT_TARGET}</string>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF
}

echo "Building one macOS app for Intel and Apple Silicon..."

for cmd in curl clang go lipo otool xcrun make perl; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "Missing required tool: $cmd"
        if [[ "$cmd" == "clang" || "$cmd" == "make" || "$cmd" == "xcrun" ]]; then
            echo "Install Xcode Command Line Tools: xcode-select --install"
        elif [[ "$cmd" == "go" ]]; then
            echo "Install Go 1.23 or newer: brew install go"
        fi
        exit 1
    fi
done

cd "$ROOT"
mkdir -p "$BUILD_DIR"

for goarch in arm64 amd64; do
    build_static_openssl "$goarch"
    build_slice "$goarch"
done

lipo -create \
    "$BUILD_DIR/ravro_dcrpt-arm64" \
    "$BUILD_DIR/ravro_dcrpt-amd64" \
    -output "$BUILD_DIR/ravro_dcrpt"
lipo -create \
    "$BUILD_DIR/ravro_dcrpt_gui-arm64" \
    "$BUILD_DIR/ravro_dcrpt_gui-amd64" \
    -output "$BUILD_DIR/ravro_dcrpt_gui"
chmod +x "$BUILD_DIR/ravro_dcrpt" "$BUILD_DIR/ravro_dcrpt_gui"

assert_static_openssl "$BUILD_DIR/ravro_dcrpt_gui" arm64
assert_static_openssl "$BUILD_DIR/ravro_dcrpt_gui" x86_64

write_app_bundle "$BUILD_DIR/ravro_dcrpt_gui"
tar -czf "$BUILD_DIR/ravro_dcrpt-darwin-universal.tar.gz" -C "$BUILD_DIR" "Ravro Decryption Tool.app"

echo
echo "Universal binary:"
lipo -info "$BUILD_DIR/ravro_dcrpt_gui"
echo
echo "Package: build/ravro_dcrpt-darwin-universal.tar.gz"
echo "This file runs on Intel Macs (x86_64) and Apple Silicon (arm64)."
echo "Run: open \"build/Ravro Decryption Tool.app\""

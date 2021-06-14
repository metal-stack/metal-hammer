#!/usr/bin/env bash
#
# Intentionally added to download metal-hammer-kernel.
set -e

KERNEL=https://github.com/metal-stack/kernel/releases/download/5.10.28-55/metal-kernel

dirty() {
    curl \
        --fail \
        --location \
        --remote-name \
        --silent \
        "${KERNEL}.md5"
    local res=$(md5sum --check "${KERNEL}.md5" 2>/dev/null 1>&2; echo $?)
    echo "${res}"
}

download() {
    curl \
        --fail \
        --location \
        --remote-name \
        "${KERNEL}"
}

download_if_dirty() {
    local isDirty=$(dirty)
    if [[ "$isDirty" = "1" ]]
    then
        echo "Downloading ..."
        download
    fi
}

download_if_dirty

# Ensure files remain writable for group to enable re-download (libvirt takes ownership for unknown reason)
chmod 660 metal-kernel 2>/dev/null || true

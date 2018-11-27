#!/usr/bin/env bash
#
# Intentionally added to download metal-hammer-kernel.
set -e

BLOBSTORE=https://blobstore.fi-ts.io

dirty() {
    curl \
        --fail \
        --location \
        --remote-name \
        --silent \
        "${BLOBSTORE}/metal/images/metal-hammer/${1}.md5"
    local res=$(md5sum --check $(basename "${1}") 2>/dev/null 1>&2; echo $?)
    echo "${res}"
}

download() {
    curl \
        --fail \
        --location \
        --remote-name \
        "${BLOBSTORE}/metal/images/metal-hammer/${1}"
}

download_if_dirty() {
    local isDirty=$(dirty "${1}")
    if [[ "$isDirty" = "1" ]]
    then
        echo "Downloading ${1}..."
        download ${1}
    fi
}

for i in "dev/metal-hammer-kernel" "metal-hammer-initrd.img.lz4"
do
    download_if_dirty $i
done

# Ensure files remain writable for group to enable re-download (libvirt takes ownership for unkown reson)
chmod 660 metal-hammer-* 2>/dev/null || true

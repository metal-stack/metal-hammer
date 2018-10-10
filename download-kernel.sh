#!/bin/bash

if [ ! -f metal-hammer-kernel ]; then
  curl -fSO https://blobstore.fi-ts.io/metal/images/metal-hammer-kernel
fi
#!/bin/bash

# clean up
echo "-> running clean up...."
rm -rf output/*

# install gox
if ! which gox > /dev/null; then
    echo "-> installing gox..."
    go get -u github.com/mitchellh/gox
fi

# build
echo "-> building..."
gox \
-os="linux" \
-arch="amd64" \
-output "output/{{.OS}}_{{.Arch}}/terrafile" \
.

# Zip and copy to the dist dir
echo ""
echo "Packaging..."
for PLATFORM in $(find ./output -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"
    
    pushd $PLATFORM >/dev/null 2>&1
    zip ../terrafile_${OSARCH}.zip ./*
done

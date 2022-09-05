#!/usr/bin/env bash
set -e

scriptDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${scriptDir}/..

for PLATFORM in $(find ./bin/dist -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"

    pushd "$PLATFORM" >/dev/null 2>&1
    if [[ ${OSARCH} = windows* ]] ; then
        zip ../waas_${OSARCH}.zip ./*
    else
        tar czvf ../waas_${OSARCH}.tgz ./*
    fi
    popd >/dev/null 2>&1
    rm -fr "$PLATFORM"
done

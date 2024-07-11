#!/bin/bash

protoBaseDir="protos"
outBaseDir="tools/gen/grpc"

if [ ! -d "$outBaseDir" ]; then
    mkdir -p "$outBaseDir"
fi

for versionDir in "$protoBaseDir"/*; do
    if [ -d "$versionDir" ]; then
        versionName=$(basename "$versionDir")
        protoSrcDir="$protoBaseDir/$versionName"
        outDir="$outBaseDir/$versionName"

        if [ ! -d "$outDir" ]; then
            mkdir -p "$outDir"
        fi

        fullProtoSrcDir=$(realpath "$protoSrcDir")

        for protoFile in "$protoSrcDir"/*.proto; do
            baseName=$(basename "$protoFile" .proto)
            targetDir="$outDir/$baseName"

            if [ ! -d "$targetDir" ]; then
                mkdir -p "$targetDir"
            fi

            protocCmd="protoc --proto_path=$fullProtoSrcDir --go_out=$targetDir --go_opt=paths=source_relative --go-grpc_out=$targetDir --go-grpc_opt=paths=source_relative $protoFile"
            echo "Running: $protocCmd"
            eval "$protocCmd"
        done
    fi
done

#!/bin/bash

protoBaseDir="protos"
outBaseDir="tools/gen/grpc"

if [ ! -d "$outBaseDir" ]; then
  mkdir -p "$outBaseDir"
fi

fullProtoBaseDir=$(realpath "$protoBaseDir")

find "$protoBaseDir" -name "*.proto" | while read -r protoFile; do
  protocCmd="protoc --proto_path=$fullProtoBaseDir --go_out=$outBaseDir --go_opt=paths=source_relative --go-grpc_out=$outBaseDir --go-grpc_opt=paths=source_relative $protoFile"
  eval "$protocCmd"
done

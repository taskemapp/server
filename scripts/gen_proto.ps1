$protoBaseDir = "protos"
$outBaseDir = "tools/gen/grpc"

if (-not (Test-Path -Path $outBaseDir)) {
    New-Item -ItemType Directory -Path $outBaseDir
}

$fullProtoBaseDir = (Resolve-Path $protoBaseDir).Path
$fullProtoBaseDir = $fullProtoBaseDir -replace '\\', '/'

Get-ChildItem -Path $protoBaseDir -Recurse -Filter *.proto | ForEach-Object {
    $protoFile = $_.FullName
    $protocCmd = "protoc --proto_path=$fullProtoBaseDir --go_out=$outBaseDir --go_opt=paths=source_relative --go-grpc_out=$outBaseDir --go-grpc_opt=paths=source_relative $protoFile"
    Invoke-Expression $protocCmd
}

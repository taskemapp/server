$protoBaseDir = "protos"
$outBaseDir = "tools/gen/grpc"

if (-not (Test-Path -Path $outBaseDir)) {
    New-Item -ItemType Directory -Path $outBaseDir
}

Get-ChildItem -Path $protoBaseDir -Directory | ForEach-Object {
    $versionDir = $_.FullName
    $versionName = $_.Name
    $protoSrcDir = Join-Path -Path $protoBaseDir -ChildPath $versionName
    $outDir = Join-Path -Path $outBaseDir -ChildPath $versionName

    if (-not (Test-Path -Path $outDir)) {
        New-Item -ItemType Directory -Path $outDir
    }

    $fullProtoSrcDir = (Resolve-Path $protoSrcDir).Path
    $fullProtoSrcDir = $fullProtoSrcDir -replace '\\', '/'

    Get-ChildItem -Path $protoSrcDir -Filter *.proto | ForEach-Object {
        $protoFile = $_.FullName
        $baseName = [System.IO.Path]::GetFileNameWithoutExtension($protoFile)
        $targetDir = Join-Path -Path $outDir -ChildPath $baseName

        if (-not (Test-Path -Path $targetDir)) {
            New-Item -ItemType Directory -Path $targetDir
        }

        $protocCmd = "protoc --proto_path=$fullProtoSrcDir --go_out=$targetDir --go_opt=paths=source_relative --go-grpc_out=$targetDir --go-grpc_opt=paths=source_relative $protoFile"
        Invoke-Expression $protocCmd
    }
}

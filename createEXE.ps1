$ErrorActionPreference = 'Stop'
$version = '1.5.22'
$exe = "mkLINK_v$version.exe"
$go = Get-Command go -ErrorAction Stop
& $go.Source env GOOS GOARCH CGO_ENABLED
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
$env:CGO_ENABLED = '0'
& $go.Source vet ./...
& $go.Source test ./...
& $go.Source build -trimpath -ldflags='-s -w -H=windowsgui' -o $exe .
Write-Host "Built $exe"

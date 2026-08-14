# Build for Linux
$env:GOOS = "linux"
$env:GOARCH = "amd64"

$deployIP = "192.168.18.162"  # GoPM-API Container

go build -o ProxmoxMgr_API

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed."
    exit 1
}

# Deploy to server - app
scp .\ProxmoxMgr_API root@[$deployIP]:/home/api/ProxmoxMgr_API

if ($LASTEXITCODE -ne 0) {
    Write-Host "Copy failed - .\ProxmoxMgr_API"
    exit 1
}

# Deploy to server - env
scp .\pm.env root@[$deployIP]:/home/api/pm.env

if ($LASTEXITCODE -ne 0) {
    Write-Host "Copy failed - \pm.env"
    exit 1
}

# CHMOD API
ssh "root@$deployIP" "chmod +x /home/api/ProxmoxMgr_API;"

if ($LASTEXITCODE -ne 0) {
    Write-Host "SSH Failed - CHMOD."
    exit 1
}

Write-Host "Deployment successful!"
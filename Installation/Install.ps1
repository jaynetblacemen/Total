$ErrorActionPreference = "Stop"

$AppName = "total"
$InstallDir = "$env:LOCALAPPDATA\Total"
$ExePath = "$InstallDir\$AppName.exe"

$DownloadUrl = "https://github.com/jaynetblacemen/Total/releases/latest/download/total-windows-amd64.exe"

Write-Host "Installing Total..." -ForegroundColor Cyan

# Create directory
if (!(Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

# Download binary
Write-Host "Downloading binary..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $ExePath

# Add to PATH (user scope)
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$UserPath;$InstallDir",
        "User"
    )
    Write-Host "Added to PATH. Restart terminal after install."
}

Write-Host "Installation complete! "
Write-Host "Run: total"

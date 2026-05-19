$ErrorActionPreference = "Stop"

$Repo = "darimuri/open-notebook-cli"
$InstallDir = "$env:LOCALAPPDATA\open-notebook\bin"
$BinaryName = "open-notebook.exe"

# Detect architecture
$Arch = $env:PROCESSOR_ARCHITECTURE
if ($Arch -eq "AMD64") { $Arch = "amd64" }
elseif ($Arch -eq "ARM64") { $Arch = "arm64" }
else { Write-Error "Unsupported architecture: $Arch" }

# Get latest version
$LatestUrl = "https://api.github.com/repos/$Repo/releases/latest"
try {
    $Latest = (Invoke-RestMethod $LatestUrl -UseBasicParsing).tag_name
} catch {
    Write-Error "Failed to get latest version: $_"
}

if (-not $Latest) {
    Write-Error "Failed to get latest version"
}

$FileName = "open-notebook-windows-$Arch.exe"
$DownloadUrl = "https://github.com/$Repo/releases/download/$Latest/$FileName"

Write-Host "Installing $Repo $Latest..."
Write-Host "Downloading $DownloadUrl..."

# Create install directory if not exists
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# Download and install
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile "$InstallDir\$BinaryName" -UseBasicParsing
} catch {
    Write-Error "Failed to download: $_"
}

# Add to PATH if not already
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
    Write-Host "Added $InstallDir to PATH"
}

Write-Host "Installed to $InstallDir\$BinaryName"

# Verify
& "$InstallDir\$BinaryName" --version
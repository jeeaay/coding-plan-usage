param(
  [string]$TargetDir = ""
)

$ErrorActionPreference = "Stop"
$repo = "jeeaay/coding-plan-usage"
$archRaw = $env:PROCESSOR_ARCHITECTURE

if ([string]::IsNullOrWhiteSpace($TargetDir)) {
  $TargetDir = $PSScriptRoot
}

switch ($archRaw.ToLower()) {
  "amd64" { $goarch = "amd64" }
  "x86_64" { $goarch = "amd64" }
  "arm64" { $goarch = "arm64" }
  default { throw "unsupported arch: $archRaw" }
}

$goos = "windows"
$asset = "coding-plan-usage-$goos-$goarch.zip"
$url = "https://github.com/$repo/releases/latest/download/$asset"

New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null
$archivePath = Join-Path $TargetDir $asset
$bundleDir = Join-Path $TargetDir "coding-plan-usage-$goos-$goarch-bundle"

Invoke-WebRequest -Uri $url -OutFile $archivePath
Expand-Archive -Path $archivePath -DestinationPath $TargetDir -Force

$binaryPath = Join-Path $bundleDir "coding-plan-usage.exe"
if (-not (Test-Path $binaryPath)) {
  throw "binary not found: $binaryPath"
}

Write-Output "installed_bundle=$bundleDir"
Write-Output "binary=$binaryPath"

@echo off
setlocal enabledelayedexpansion

set "TargetDir=%~1"
set "repo=jeeaay/coding-plan-usage"
set "goos=windows"

if "%TargetDir%"=="" (
    set "TargetDir=%~dp0"
)

set "goarch="
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" set "goarch=amd64"
if "%PROCESSOR_ARCHITECTURE%"=="x86_64" set "goarch=amd64"
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" set "goarch=arm64"

if "!goarch!"=="" (
    echo unsupported arch: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)

set "asset=coding-plan-usage-%goos%-%goarch%.zip"
set "url=https://github.com/%repo%/releases/latest/download/%asset%"

if not exist "!TargetDir!" mkdir "!TargetDir!"
set "archivePath=!TargetDir!\%asset%"
set "bundleDir=!TargetDir!\coding-plan-usage-%goos%-%goarch%-bundle"

echo Downloading %url%...
powershell -Command "Invoke-WebRequest -Uri '%url%' -OutFile '%archivePath%'"

echo Extracting...
powershell -Command "Expand-Archive -Path '%archivePath%' -DestinationPath '%TargetDir%' -Force"

set "binaryPath=!bundleDir!\coding-plan-usage.exe"
if not exist "!binaryPath!" (
    echo binary not found: !binaryPath!
    exit /b 1
)

echo installed_bundle=!bundleDir!
echo binary=!binaryPath!

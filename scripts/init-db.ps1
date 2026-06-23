$ErrorActionPreference = "Stop"

$backendPath = Join-Path $PSScriptRoot "..\backend"
Push-Location $backendPath
try {
  & go run ./cmd/cli db init @args
  if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
  }
} finally {
  Pop-Location
}

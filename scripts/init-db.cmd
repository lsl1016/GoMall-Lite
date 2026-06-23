@echo off
setlocal
pushd "%~dp0..\backend" || exit /b 1
go run ./cmd/cli db init %*
set EXIT_CODE=%ERRORLEVEL%
popd
exit /b %EXIT_CODE%

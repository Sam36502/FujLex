@ECHO off

ECHO Building Windows binary...
go env -w GO111MODULE=on
go build -o bin\fujlex.exe src\main.go
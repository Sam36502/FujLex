@ECHO off

ECHO Building Windows binary...
go env -w GO111MODULE=on
go build -o bin\win\4RCH.exe src\main.go
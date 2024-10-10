#!/bin/sh
go build server.go
while true; do ./server; done
#!/bin/sh
cd cli && GOOS=linux GOARCH=amd64 go build
cd ../dbserver && GOOS=linux GOARCH=amd64 go build
cd ../frontend && GOOS=linux GOARCH=amd64 go build

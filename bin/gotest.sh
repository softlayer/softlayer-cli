#!/usr/bin/bash
set -e

go vet $(go list ./... | grep -v "fixtures" | grep -v "vendor")
go test $(go list ./... | grep -v "fixtures" | grep -v "vendor")
gosec -exclude-dir=fixture -exclude-dir=plugin/resources -quiet ./...
go build
go vet $(go list ./... | grep -v "fixtures" | grep -v "vendor")
go test $(go list ./... | grep -v "fixtures" | grep -v "vendor")
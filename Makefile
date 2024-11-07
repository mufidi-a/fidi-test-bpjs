run-test-services:
	go test -covermode=count -coverpkg=./service/...,./controller/... -coverprofile cover.out -v ./controller ./service/...

html-run-test-services:
	cd service && go tool cover -html cover.out -o cover.html

start:
	go run cmd/main.go

mock-repo:
	cd repository && mockery --all --case=underscore && cd ../

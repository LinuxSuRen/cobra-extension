build: lint fmt
	go build

lint:
	golint ./...

fmt:
	gofmt -s -w .

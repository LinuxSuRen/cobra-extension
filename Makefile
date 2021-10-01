build: lint fmt
	go mod tidy
	go build

copy: build
	cp cobra-extension /usr/local/bin

lint:
	golint ./...

fmt:
	gofmt -s -w .

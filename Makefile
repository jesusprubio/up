.PHONY: all dev build clean dep fmt vet test

all: dev

dev:
	go run .

build:
	go build -o dist/up .

clean:
	go clean

dep:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/segmentio/golines@latest

fmt:
	golines -m 80 --shorten-comments -w .

vet:
	go vet ./...
	staticcheck ./...
	gosec -fmt=golint -quiet ./...

test:
	go test ./...

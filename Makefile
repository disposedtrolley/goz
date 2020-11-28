ENTRYPOINT:=./cmd/goz/main.go
BINARY_NAME:=goz

build:
	go build -mod=vendor -o $(BINARY_NAME) $(ENTRYPOINT)

run:
	go run -mod=vendor $(ENTRYPOINT)

.PHONY: test
test:
	go test -race ./...

.PHONY: clean test build zip

default: build

clean:
	-rm -rf bin/*

test:
	go test -race -v ./...

build:
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/handle ./cmd/handle
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/process ./cmd/process
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/save ./cmd/save
	
zip:
	@cd ./bin && zip handle.zip handle
	@cd ./bin && zip process.zip process
	@cd ./bin && zip save.zip save
	
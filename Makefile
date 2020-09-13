.PHONY: clean test build zip

default: clean build zip

clean:
	-rm -rf bin/*

test:
	go test -race -v ./...

build:
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/handler ./cmd/handler
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/processor ./cmd/processor
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/checker ./cmd/checker
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/creator ./cmd/creator
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/updater ./cmd/updater
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/dbputter ./cmd/dbputter
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/dbupdater ./cmd/dbupdater
	
zip:
	@cd ./bin && zip handler.zip handler
	@cd ./bin && zip processor.zip processor
	@cd ./bin && zip checker.zip checker
	@cd ./bin && zip creator.zip creator
	@cd ./bin && zip updater.zip updater
	@cd ./bin && zip dbputter.zip dbputter
	@cd ./bin && zip dbupdater.zip dbupdater
	
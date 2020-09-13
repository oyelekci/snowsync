.PHONY: clean test build zip

default: clean test build zip

clean:
	-rm -rf bin/*

test:
	go test -cover -v ./...

build:
	 GOOS=linux GOARCH=amd64 go build -v -o ./bin/ ./cmd/...
	
zip:
	@cd ./bin && find . -type f -exec zip -D '{}.zip' '{}' \;

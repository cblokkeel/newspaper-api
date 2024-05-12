build:
	@go build -C cmd/newspaper -o ../../bin/api

run: build
	@./bin/api

test:
	@go test -v ./...

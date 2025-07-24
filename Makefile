build:
	@go build -o bin/Bank
run: build
	@./bin/Bank
test: 
	@go test -v ./...

build: 
	@go build -o bin/api

run: build
	@./bin/api

dev: 
	air


test:
	@go test -v ./...

test-cover:
	@go test -v ./... -cover

seed:
	@go run scripts/seed.go
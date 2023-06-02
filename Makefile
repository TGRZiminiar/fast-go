# mingw32-make run
build:
	@go build -o bin/based

run: build
	@./bin/based

test:
	go test -v ./...
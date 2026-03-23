BINARY=sailo
VERSION?=dev
LDFLAGS=-ldflags "-X github.com/agawish/sailo/cmd/sailo/commands.version=$(VERSION)"

.PHONY: build test lint install clean

build:
	go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/sailo

test:
	go test ./... -v -race

lint:
	go vet ./...
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"

install: build
	cp bin/$(BINARY) $(GOPATH)/bin/$(BINARY) 2>/dev/null || cp bin/$(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -rf bin/

run: build
	./bin/$(BINARY) $(ARGS)

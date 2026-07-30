BINARY=neatcli
GOBUILD=go build
GOTEST=go test
GOFMT=go fmt
GOVET=go vet

.PHONY: build clean test fmt vet lint install run help

build:
	$(GOBUILD) -o $(BINARY) .

clean:
	rm -f $(BINARY)

test:
	$(GOTEST) ./...

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

lint: fmt vet

install: build
	mv $(BINARY) $(GOPATH)/bin/$(BINARY)

run: build
	./$(BINARY)

help:
	@echo "Targets:"
	@echo "  build    - compile the binary"
	@echo "  clean    - remove build artifacts"
	@echo "  test     - run tests"
	@echo "  fmt      - format Go source files"
	@echo "  vet      - run go vet"
	@echo "  lint     - fmt + vet"
	@echo "  install  - build and copy to GOPATH/bin"
	@echo "  run      - build and execute"

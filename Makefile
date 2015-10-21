PREFIX?=/build

GOFILES = $(shell find . -type f -name '*.go')
nginxbeat: $(GOFILES)
	go build

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm nginxbeat || true


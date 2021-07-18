.PHONY: all
all: ./bin/bot

bin/%: $(shell find . -type f -name '*.go')
	@mkdir -p $(dir $@)
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -o $@ ./cmd/$(@F)

.PHONY: clean
clean:
	rm -rf bin

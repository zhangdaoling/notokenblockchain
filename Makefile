all: build

tag := `git describe --always`
commit := `git describe --always`
buildtime := `date +%FT%T%z`
MODULE = "huobi.com/blockchain/sync"

LD_FLAGS := -ldflags "-X main.commit=${commit} -X main.buildTime=${buildtime}"
CURRENT_DIR=$(shell pwd)

GOFMT := $(shell command -v gofmt 2> /dev/null)
GOIMPORTS := $(shell command -v goimports 2> /dev/null)
GOLINT :=  $(shell command -v golint 2> /dev/null)
CILINT := $(shell command -v golangci-lint 2> /dev/null)

style: clean
ifndef GOFMT
	$(error "gofmt is not available please install gofmt")
endif
ifndef GOIMPORTS
	$(error "goimports is not available please install goimports")
endif
	@echo ">> checking code style"
	@! find . -path ./vendor -prune -o -name '*.go' -print | xargs gofmt -d | grep '^'
	@! find . -path ./vendor -prune -o -name '*.go' -print | xargs goimports -d -local ${MODULE} | grep '^'

format:
ifndef GOFMT
	$(error "gofmt is not available please install gofmt")
endif
ifndef GOIMPORTS
	$(error "goimports is not available please install goimports")
endif
	@go fmt ./...
	@find . -path ./vendor -prune -o -name '*.go' -print | xargs goimports -l -local ${MODULE} | xargs goimports -l -local ${MODULE} -w

lint:
ifndef GOLINT
	$(error "golint is not available please install golint")
endif
	@echo ">> checking code lint"
	@! go list ./... | xargs golint | sed "s:^$(CURRENT_DIR)/::" | grep '^'

cilint:
ifndef CILINT
	$(error "golangci-lint is not available please install golangci-lint")
endif
	@golangci-lint run

test: style lint
	@echo ">> testing"
	@go test -cover ./...

wocao: test
	@go build ${LD_FLAGS} -o build/chain app/main.go

vendor:
	go mod vendor

build: wocao

clean:
	@rm -rf build

cleanAll: clean
	@rm -rf vendor

.PHONY: sync server style format test build clean cleanAll vendor

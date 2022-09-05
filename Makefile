GO_TAGS=
GO_ARGS=-tags '$(GO_TAGS)'
VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
COMMIT := $(shell git rev-parse --short HEAD)

LDFLAGS := $(LDFLAGS) -X main.commit=$(COMMIT)
ifdef VERSION
	LDFLAGS += -X main.version=$(VERSION)
endif

export XC_ARCH=386 amd64 arm arm64
export XC_OS=solaris darwin freebsd linux windows
export CGO_ENABLED=0
export GO_BUILD=env GO111MODULE=on go build $(GO_ARGS) -ldflags "$(LDFLAGS)"
export GOX_BUILD=env GO111MODULE=on gox -os="$(XC_OS)" -arch="$(XC_ARCH)" -osarch="!darwin/arm !darwin/arm64 !darwin/386" \
                -output "bin/dist/{{.OS}}_{{.Arch}}/{{.Dir}}" \
                $(GO_ARGS) -ldflags "$(LDFLAGS)"
export GO_TEST=env GOTRACEBACK=all GO111MODULE=on go test $(GO_ARGS)
export GO_TEST_IT=env GOTRACEBACK=all GO111MODULE=on go test -count=1 -tags=integration -v
export GO_VET=env GO111MODULE=on go vet $(GO_ARGS)
export PATH := $(PWD)/bin/:$(PATH)


# All go source files
GO_SOURCES := $(shell find . -name '*.go')
SOURCES := $(GO_SOURCES) go.mod go.sum


# List of binary cmds to build
CMDS := \
	bin/hpcwaas-api \
	bin/waas

all: fmt $(CMDS)

#
# Define targets for commands
#
$(CMDS):  fmt $(SOURCES)
	$(GO_BUILD) -o $@ ./cmd/$(shell basename "$@")


# Ease of use build for just the go binary
hpcwaas-api: bin/hpcwaas-api
waas: bin/waas

dist: fmt
	$(GOX_BUILD) ./cmd/waas
	./scripts/dist.sh

tools:
	env GO111MODULE=on go get -tags tools

fmt: generate
	@gofmt -w -s $(GO_SOURCES)

tidy:
	GO111MODULE=on go mod tidy

test-go:
	$(GO_TEST) ./...

test: test-go


test-go-race:
	CGO_ENABLED=1 $(GO_TEST) -v -race -count=1 ./...

test-json:
	./bin/gotestsum --jsonfile tests-reports.json  -- -tags "$(BUILD_TAGS)" $(TESTARGS) -count=1 -p 1 ./...

vet:
	$(GO_VET) -v ./...

bench:
	$(GO_TEST) -bench=. -run=^$$ ./...

build: all

clean:
	$(RM) -r bin

generate: tools
#	@go install github.com/abice/go-enum
	@go generate ./...

# .PHONY targets represent actions that do not create an actual file.
.PHONY: all fmt tidy generate test test-go test-go-race bench clean vet

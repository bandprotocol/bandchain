PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin

ifeq ($(LEDGER_ENABLED),true)
	build_tags += ledger
endif

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=bandchain \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=bandd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=bandcli \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandcli
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandoracled2

faucet: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/faucet

release: go.sum
	env GOOS=linux GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandd_linux_amd64 $(BUILD_FLAGS) ./cmd/bandd
	env GOOS=darwin GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandd_darwin_amd64 $(BUILD_FLAGS) ./cmd/bandd
	env GOOS=windows GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandd_windows_amd64 $(BUILD_FLAGS) ./cmd/bandd
	env GOOS=linux GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandcli_linux_amd64 $(BUILD_FLAGS) ./cmd/bandcli
	env GOOS=darwin GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandcli_darwin_amd64 $(BUILD_FLAGS) ./cmd/bandcli
	env GOOS=windows GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandcli_windows_amd64 $(BUILD_FLAGS) ./cmd/bandcli
	env GOOS=linux GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandoracled_linux_amd64 $(BUILD_FLAGS) ./cmd/bandoracled2
	env GOOS=darwin GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandoracled_darwin_amd64 $(BUILD_FLAGS) ./cmd/bandoracled2
	env GOOS=windows GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandoracled_windows_amd64 $(BUILD_FLAGS) ./cmd/bandoracled2

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

###############################################################################
###                                Protobuf                                 ###
###############################################################################
PREFIX ?= /usr/local
PROTOC_VERSION ?= 3.11.2
PROTOC_ZIP ?= protoc-3.11.2-linux-x86_64.zip
# ifeq ($(UNAME_S),Linux)
#   PROTOC_ZIP ?= protoc-3.11.2-linux-x86_64.zip
# endif
# ifeq ($(UNAME_S),Darwin)
#   PROTOC_ZIP ?= protoc-3.11.2-osx-x86_64.zip
# endif

proto-all: proto-gen

proto-gen:
	@./scripts/protocgen.sh

proto-tools-stamp:
	@echo "Installing protoc compiler..."
	@echo  "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"
	@(cd /tmp; \
	curl -OL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"; \
	unzip -o ${PROTOC_ZIP} -d $(PREFIX) bin/protoc; \
	unzip -o ${PROTOC_ZIP} -d $(PREFIX) 'include/*'; \
	rm -f ${PROTOC_ZIP})

	@echo "Installing protoc-gen-buf-check-breaking..."
	@curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/protoc-gen-buf-check-breaking-${UNAME_S}-${UNAME_M}" \
    -o "${BIN}/protoc-gen-buf-check-breaking" && \
	chmod +x "${BIN}/protoc-gen-buf-check-breaking"

	@echo "Installing buf..."
	@curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-${UNAME_S}-${UNAME_M}" \
    -o "${BIN}/buf" && \
	chmod +x "${BIN}/buf"

	touch $@

protoc-gen-gocosmos:
	@echo "Installing protoc-gen-gocosmos..."
	@go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos

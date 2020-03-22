PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BINDIR ?= $(GOPATH)/bin

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=d3n \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=bandd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=bandcli \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: lint install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandcli
	# go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandsv
	# go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandoracled

release: go.sum
	env GOOS=linux GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandd_linux_amd64 $(BUILD_FLAGS) ./cmd/bandd
	env GOOS=linux GOARCH=arm \
		go build -mod=readonly -o ./build/bandd_linux_arm $(BUILD_FLAGS) ./cmd/bandd
	env GOOS=darwin GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandd_darwin_amd64 $(BUILD_FLAGS) ./cmd/bandd
	env GOOS=windows GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandd_windows_amd64 $(BUILD_FLAGS) ./cmd/bandd
	env GOOS=linux GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandcli_linux_amd64 $(BUILD_FLAGS) ./cmd/bandcli
	env GOOS=linux GOARCH=arm \
		go build -mod=readonly -o ./build/bandcli_linux_arm $(BUILD_FLAGS) ./cmd/bandcli
	env GOOS=darwin GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandcli_darwin_amd64 $(BUILD_FLAGS) ./cmd/bandcli
	env GOOS=windows GOARCH=amd64 \
		go build -mod=readonly -o ./build/bandcli_windows_amd64 $(BUILD_FLAGS) ./cmd/bandcli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

update-swagger-docs:
	$(BINDIR)/statik -src=client/lcd/swagger-ui -dest=client/lcd -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
			echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
			exit 1;\
	else \
		echo "\033[92mSwagger docs are in sync\033[0m";\
	fi

lint:
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

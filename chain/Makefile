PACKAGES=$(shell go list ./... | grep -v '/simulation')

COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=zoracle \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=bandd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=bandcli \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: lint install

install: go.sum
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandd
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandcli
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandsv
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/bandoracled

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

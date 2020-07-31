.PHONY: docker-images docker-image-linux docker-image-osx

DOCKER_TAG := 0.0.1
USER_ID := $(shell id -u)
USER_GROUP = $(shell id -g)

docker-image-linux:
	DOCKER_BUILDKIT=1 docker build .. -t owasm/go-ext-builder:$(DOCKER_TAG)-linux -f Dockerfile.linux

docker-image-osx:
	DOCKER_BUILDKIT=1 docker build .. -t owasm/go-ext-builder:$(DOCKER_TAG)-osx -f Dockerfile.osx

docker-images: docker-image-linux docker-image-osx

# and use them to compile release builds
release:
	rm -rf target/release
	docker run --rm -u $(USER_ID):$(USER_GROUP) -v $(shell pwd)/..:/code owasm/go-ext-builder:$(DOCKER_TAG)-osx
	rm -rf target/release
	docker run --rm -u $(USER_ID):$(USER_GROUP) -v $(shell pwd)/..:/code owasm/go-ext-builder:$(DOCKER_TAG)-linux

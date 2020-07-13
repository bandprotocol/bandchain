.PHONY: lambda

VERSION=1.0.1
CURRENT_DIR = $(shell pwd)

lambda:
	docker run --rm -v $(CURRENT_DIR)/lambda:/pack amazonlinux:latest ./pack/pack.sh

docker-build:
	docker build docker -t bandprotocol/runtime:$(VERSION)

docker-push: docker-build
	docker push bandprotocol/runtime:$(VERSION)

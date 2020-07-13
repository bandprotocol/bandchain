.PHONY: lambda

VERSION=1.0.1
CURRENT_DIR=$(shell pwd)

lambda:
	cp $(CURRENT_DIR)/requirements.txt $(CURRENT_DIR)/lambda
	docker run --rm -v $(CURRENT_DIR)/lambda:/pack amazonlinux:latest ./pack/pack.sh
	rm $(CURRENT_DIR)/lambda/requirements.txt

docker-build:
	docker build docker -t bandprotocol/runtime:$(VERSION)

docker-push: docker-build
	docker push bandprotocol/runtime:$(VERSION)

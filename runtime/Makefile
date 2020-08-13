.PHONY: lambda

VERSION=1.0.3
CURRENT_DIR=$(shell pwd)
RUNTIME_VERSION="lambda:v1.0.3"

lambda:
	cp $(CURRENT_DIR)/lambda/lambda_function.py $(CURRENT_DIR)/lambda/lambda_function.bak
	RUNTIME_VERSION=$(RUNTIME_VERSION) envsubst < $(CURRENT_DIR)/lambda/lambda_function.bak > $(CURRENT_DIR)/lambda/lambda_function.py
	cp $(CURRENT_DIR)/requirements.txt $(CURRENT_DIR)/lambda
	docker run --rm -v $(CURRENT_DIR)/lambda:/pack amazonlinux:latest ./pack/pack.sh
	rm $(CURRENT_DIR)/lambda/requirements.txt
	mv $(CURRENT_DIR)/lambda/lambda_function.bak $(CURRENT_DIR)/lambda/lambda_function.py

docker-build:
	cp $(CURRENT_DIR)/requirements.txt $(CURRENT_DIR)/docker
	docker build docker -t bandprotocol/runtime:$(VERSION)
	rm $(CURRENT_DIR)/docker/requirements.txt

docker-push: docker-build
	docker push bandprotocol/runtime:$(VERSION)

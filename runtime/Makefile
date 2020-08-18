.PHONY: lambda

VERSION=1.1.1
CURRENT_DIR=$(shell pwd)

lambda:
	cp $(CURRENT_DIR)/lambda/lambda_function.py $(CURRENT_DIR)/lambda/lambda_function.bak
	RUNTIME_VERSION="lambda:$(VERSION)" envsubst < $(CURRENT_DIR)/lambda/lambda_function.bak > $(CURRENT_DIR)/lambda/lambda_function.py
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

gcf:
	cp $(CURRENT_DIR)/google-cloud-functions/google_cloud_function.py $(CURRENT_DIR)/google-cloud-functions/google_cloud_function.bak
	RUNTIME_VERSION="google-cloud-function:$(VERSION)" envsubst < $(CURRENT_DIR)/google-cloud-functions/google_cloud_function.bak > $(CURRENT_DIR)/google-cloud-functions/google_cloud_function.py
	cd google-cloud-functions && zip -g google_cloud_function.zip google_cloud_function.py requirements.txt & cd ..
	mv $(CURRENT_DIR)/google-cloud-functions/google_cloud_function.bak $(CURRENT_DIR)/google-cloud-functions/google_cloud_function.py



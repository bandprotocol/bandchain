VERSION=1.0.1

pack-lambda:
	cd lambda && ./pack.sh

docker-build:
	docker build docker -t bandprotocol/runtime:$(VERSION)

docker-push: docker-build
	docker push bandprotocol/runtime:$(VERSION)

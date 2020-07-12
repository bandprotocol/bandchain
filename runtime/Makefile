VERSION=1.0.1

docker-build:
	docker build docker -t bandprotocol/runtime:$(VERSION)

docker-push: docker-build
	docker push bandprotocol/runtime:$(VERSION)

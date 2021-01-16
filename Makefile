.PHONY: all

all: build test
IMG="kaizheh/process-connector"
VERSION=$(shell cat version)

test:
	@echo "+ $@"
	go test ./...
build:
	@echo "+ $@"
	./scripts/build
build-image:
	@echo "+ $@"
	docker build -f container/Dockerfile -t ${IMG}:${VERSION} .
push-image:
	@echo "+ $@"
	docker push ${IMG}:${VERSION}
	docker tag ${IMG}:${VERSION} ${IMG}:latest
	docker push ${IMG}:latest
docker-run:
	@echo "+ $@"
	docker run --rm -v /proc:/host/proc:ro -v /var/lib/docker/containers:/host/containers --net host --cap-add net_admin kaizheh/process-connector:${VERSION}

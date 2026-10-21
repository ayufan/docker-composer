export REGISTRY ?= ayufan/docker-composer
export TAG ?= beta
export BUILD_ARCHS ?= amd64 arm32v7 arm64v8
HOST_PLATFORM ?= $(shell docker version --format '{{.Server.Os}}/{{.Server.Arch}}')

all: $(addsuffix -docker-build, $(BUILD_ARCHS))

%-docker-build:
	docker build \
		--tag $(REGISTRY):$(TAG)-$* \
		--build-arg ARCH=linux/$(subst arm32v,arm/v,$(subst arm64v,arm64/v,$*)) \
		--build-arg REPO=$*/ \
		-f Dockerfile \
		.

%-dockerhub: %-docker-build
	docker push $(REGISTRY):$(TAG)-$*

dockerhub: $(addsuffix -dockerhub, $(BUILD_ARCHS))
	-docker manifest rm $(REGISTRY):$(TAG)
	docker manifest create $(REGISTRY):$(TAG) \
		$(addprefix $(REGISTRY):$(TAG)-, $(BUILD_ARCHS))
	docker manifest push $(REGISTRY):$(TAG)

tag:
	+make dockerhub TAG=$(shell git describe --exact-match)

tag-latest:
	+make dockerhub TAG=latest

dev-env:
	docker run --rm -it -w /go/src/github.com/ayufan/docker-composer -v $(CURDIR):/go/src/github.com/ayufan/docker-composer golang:alpine /bin/sh

version:
	docker build \
		--tag $(REGISTRY):$(TAG)-version \
		--build-arg ARCH=$(HOST_PLATFORM) \
		-f Dockerfile \
		.
	docker run --rm --entrypoint /bin/sh $(REGISTRY):$(TAG)-version -c '\
		docker --version; \
		docker-compose --version; \
		git --version'

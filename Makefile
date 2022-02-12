export REGISTRY ?= ayufan/docker-composer
export TAG ?= beta
export BUILD_ARCHS ?= amd64 arm32v7 arm64v8

all: $(addsuffix -docker-build, $(BUILD_ARCHS))

%-docker-build:
	docker build \
		--tag $(REGISTRY):$(TAG)-$* \
		--build-arg ARCH=$*/ \
		-f Dockerfile \
		.

%-dockerhub: %-docker-build
	docker push $(REGISTRY):$(TAG)-$*

dockerhub: $(addsuffix -dockerhub, $(BUILD_ARCHS))
	-rm -rf ~/.docker/manifests
	docker manifest create $(REGISTRY):$(TAG) \
		$(addprefix $(REGISTRY):$(TAG)-, $(BUILD_ARCHS))
	docker manifest push $(REGISTRY):$(TAG)

tag:
	make dockerhub TAG=$(shell git describe --exact-match)

tag-latest:
	make dockerhub TAG=latest


RELEASE_TAG:=$(shell git describe --tags --exact-match --match 'v[0-9]*.[0-9]*.[0-9]*' HEAD 2>/dev/null)
TAG ?= $(or $(patsubst v%,%,$(RELEASE_TAG)),dev)
NAME := dgeorgievski-hello
MODULE:=$(shell awk '$$1 == "module" { print $$2 }' go.mod)
DOCKER_IMAGE_NAME ?= ghcr.io/dgeorgievski/hello
GIT_COMMIT:=$(shell git rev-parse --short HEAD)
GIT_DATE:=$(shell git show -s --format=%cI HEAD)
VERSION:=$(shell grep 'var VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')
EXTRA_RUN_ARGS?=
BUILD_OUTPUT ?= ./bin/$(NAME)

# docker push registry.gitlab.casa-systems.com/dimitar.georgievski/5g-metrics
run:
	go run main.go  $(EXTRA_RUN_ARGS)

.PHONY: test
test:
	go test ./... -coverprofile cover.out

.PHONY: build check-build
build:
	go build \
		-ldflags "-s -w -X $(MODULE)/pkg/version.VERSION=$(TAG) -X $(MODULE)/pkg/version.COMMIT=$(GIT_COMMIT) -X $(MODULE)/pkg/version.DATE=$(GIT_DATE)" \
		-o "$(BUILD_OUTPUT)" .

check-build: build
	sh scripts/check-build.sh "$(BUILD_OUTPUT)" "$(TAG)" "$(GIT_COMMIT)" "$(GIT_DATE)"

tidy:
	go mod tidy

fmt:
	gofmt -l -s -w ./

docker-build:
	docker build \
	--build-arg VERSION="$(TAG)" \
	--build-arg COMMIT="$(GIT_COMMIT)" \
	--build-arg BUILD_DATE="$(GIT_DATE)" \
	-f docker/Dockerfile \
	 -t $(DOCKER_IMAGE_NAME):$(TAG) .

docker-test:
	@echo "Coming soon"

docker-push:
	# docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	@echo "Coming soon"

kind-docker-upload: docker-build
	 kind --name hello load docker-image $(DOCKER_IMAGE_NAME):$(TAG)

version-set:
	# @next="$(TAG)" && \
	# current="$(VERSION)" && \
	# /usr/bin/sed -i "s/$$current/$$next/g" pkg/version/version.go && \
	# /usr/bin/sed -i "s/hello:$$current/hello:$$next/g" deploy/manifests/deployment.yaml && \
	# echo "Version $$next set in code and manifests"
	@echo "Coming soon"
	
version-get:
	@echo "Current version $(VERSION)"

tag:
	git tag v$(VERSION)
	git push origin v$(VERSION)

get-env:
	@echo "TAG: $(TAG)"
	@echo "NAME: $(NAME)" 
	@echo "DOCKER_IMAGE_NAME: $(DOCKER_IMAGE_NAME)"
	@echo "GIT_COMMIT: $(GIT_COMMIT)"
	@echo "VERSION: $(VERSION)"
	@echo "EXTRA_RUN_ARGS: $(EXTRA_RUN_ARGS)"


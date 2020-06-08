.DEFAULT_GOAL := build-all

DOCKER_NAME = golang
VERSION = 1.14.3
DOCKER_IMAGE = $(DOCKER_NAME):$(VERSION)

BUILDER_CURRENT_DIR = $(shell pwd)
BUILDER_PROJECT_DIR = $(shell git rev-parse --show-toplevel)
BUILDER_USER = --user=$(shell id -u):$(shell id -g)
BUILDER_VOLUMES = -v ${BUILDER_PROJECT_DIR}:${BUILDER_PROJECT_DIR}
BUILDER_ARGS = ${BUILDER_USER} -w ${BUILDER_CURRENT_DIR} ${BUILDER_VOLUMES} 

init:
	docker pull $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	docker tag $(DOCKER_REGISTRY)/$(DOCKER_IMAGE) $(DOCKER_IMAGE)

build-server: 
	docker run --rm -t ${BUILDER_ARGS} ${DOCKER_IMAGE} make -C cmd/server

build-client:
	docker run --rm -t ${BUILDER_ARGS} ${DOCKER_IMAGE} make -C cmd/client

build-all: build-server build-client

clean:
	make -C cmd/server clean
	make -C cmd/client clean
	rm -rf build

.PHONY: init build-server build-client build-all clean


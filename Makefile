.DEFAULT_GOAL := build-all

APP_VERSION = 1.0

build-server: 
	make -C cmd/server

build-client:
	make -C cmd/client

build-all: build-server build-client

docker:
	@cp dist/usr/bin/todoserver deploy/todoserver
	@cp dist/usr/bin/todoserverconf.json deploy/todoserverconf.json
	@cd deploy; docker build -t todoapp:${APP_VERSION} .

clean:
	make -C cmd/server clean
	make -C cmd/client clean
	rm -rf build

.PHONY: init build-server build-client build-all clean


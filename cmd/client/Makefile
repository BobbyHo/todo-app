.DEFAULT_GOAL := build

dirs := internal
search_dirs := $(addprefix ./,$(addsuffix /...,$(dirs)))

BIN_DIR = ../../dist/usr/bin

init-dist-dir:
	mkdir -p $(BIN_DIR)

build: init-dist-dir
	go clean; go build -o todoclient
	mv -f todoclient $(BIN_DIR)

rebuild: clean-artifact build

clean-artifact:
	go clean --modcache
	rm -rf $(BIN_DIR)/todoclient
	
clean: clean-artifact

.PHONY: init-dist-dir build rebuild clean-artifact clean
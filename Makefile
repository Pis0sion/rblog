APP_NAME = rblog
BUILD_DIR = $(PWD)/build
GO_BUILD = /usr/local/opt/go/libexec/bin/go build
MAIN = $(PWD)/cmd/rblog
.BUILD_GOAL := all

.PHONY: all
all: clean build run

.PHONY: clean test build run

clean:
	@rm -fr $(BUILD_DIR)

build:
	@$(GO_BUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

run:
	@$(BUILD_DIR)/$(APP_NAME)
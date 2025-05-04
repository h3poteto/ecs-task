.PHONY: all mod build

OUTPUT = ecs-task
OUTDIR = bin
BUILD_CMD = go build -a -tags netgo -installsuffix netgo -ldflags \
" \
  -extldflags '-static' \
  -X github.com/h3poteto/ecs-task/cmd.version=$(shell git describe --tag --abbrev=0) \
  -X github.com/h3poteto/ecs-task/cmd.revision=$(shell git rev-list -1 HEAD) \
  -X github.com/h3poteto/ecs-task/cmd.build=$(shell git describe --tags) \
"
ARCHITECTURE = amd64 arm64
VERSION = $(shell git describe --tag --abbrev=0)

all: mac linux windows

mod: go.mod
	go mod download

build: mod
	GOOS=linux GOARCH=amd64 $(BUILD_CMD) -o $(OUTPUT)

mac: mod
	@for arch in $(ARCHITECTURE); do \
	    GOOS=darwin GOARCH=$$arch CGO_ENABLED=0 $(BUILD_CMD) -o $(OUTPUT); \
		zip $(OUTDIR)/$(OUTPUT)_${VERSION}_darwin_$$arch.zip $(OUTPUT); \
	done

linux: mod
	@for arch in $(ARCHITECTURE); do \
		GOOS=linux GOARCH=$$arch CGO_ENABLED=0 $(BUILD_CMD) -o $(OUTPUT); \
		zip $(OUTDIR)/$(OUTPUT)_${VERSION}_linux_$$arch.zip $(OUTPUT); \
	done

windows: mod
	@for arch in $(ARCHITECTURE); do \
		GOOS=windows GOARCH=$$arch CGO_ENABLED=0 $(BUILD_CMD) -o $(OUTPUT).exe; \
		zip $(OUTDIR)/$(OUTPUT)_${VERSION}_windows_$$arch.zip $(OUTPUT).exe; \
	done

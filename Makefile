CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
CURRENT_COMMIT := $(shell git rev-parse --short HEAD)
TAG := $(shell git describe --exact-match --tags $(CURRENT_COMMIT) 2>/dev/null)

ifeq ($(TAG),)
	VERSION := $(CURRENT_BRANCH)-$(CURRENT_COMMIT)
else
	VERSION := $(TAG)
endif

.PHONY: build
build:
	@echo "Building version $(VERSION)"
	go build -ldflags="-X 'main.version=$(VERSION)'" ./cmd/kubectl-gpt

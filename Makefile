BINARY_NAME=guildops
VERSION=$(shell git describe --tags --always)

## ----------------------
## Available make targets
## ----------------------
##

default: help

help: ## Display this message
	@grep -E '(^[a-zA-Z0-9_\-\.]+:.*?##.*$$)|(^##)' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

##
## ----------------------
## Builds
## ----------------------
##

artifact: ## Compile app from sources (linux)
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ${BINARY_NAME} ./cmd/${BINARY_NAME}

artifact.osx: ## Compile app from sources (osx)
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ${BINARY_NAME} ./cmd/${BINARY_NAME}

image-ci: ## Build an image for CI Test Helm
	docker build . --tag "ghcr.io/antony-ramos/${BINARY_NAME}:ci"

##
## ----------------------
## Q.A
## ----------------------
##

qa: lint test.unit ## Run all Q.A

lint.install: ## Install Go linter
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: ## Lint source code
	golangci-lint run -v

lint.fix: ## Lint and fix source code
	golangci-lint run --fix -v

test: ## Run unit tests (go test)
	@echo "### Running unit tests ..."
	go test -v ./... -coverprofile coverage.output

##
## ----------------------
## Development
## ----------------------
##

install: ## Install required goland dependencies
	go mod download

start: ## Start project
	go run ./cmd/${BINARY_NAME}

##
## ----------------------
## Debian package
## ----------------------
##

dpkg.build: artifact ## Build debian package
	@echo "### Building debian package ..."
	@mkdir -p ./install/${BINARY_NAME}/usr/bin
	@cp ${BINARY_NAME} ./install/${BINARY_NAME}/usr/bin/${BINARY_NAME}
	@mkdir -p ./install/${BINARY_NAME}/etc/${BINARY_NAME}/config
	@cp config/config.yml ./install/${BINARY_NAME}/etc/${BINARY_NAME}/config/config.yml
	@dpkg-deb --build ./install/${BINARY_NAME} ./${BINARY_NAME}_${VERSION}.deb
	@echo "### Debian package built successfully"

dpkg.install:
	@echo "### Installing debian package ..."
	@sudo dpkg -i ./${BINARY_NAME}_${VERSION}.deb
	@echo "### Debian package installed successfully"

dpkg.remove:
	@echo "### Removing debian package ..."
	@sudo dpkg -r ${BINARY_NAME}
	@echo "### Debian package removed successfully"

##
## ----------------------
## SSH Install (Debian)
## ----------------------
##

# SSH configuration (optional)
SSH_USER ?= root
SSH_HOST ?= localhost
SSH_TARGET_DIR ?= /tmp

DEB_PACKAGE = ./${BINARY_NAME}_${VERSION}.deb

ssh.deploy: dpkg.build
	@echo "### Deploying debian package to remote machine ..."
	scp -P 22 $(DEB_PACKAGE) $(SSH_USER)@$(SSH_HOST):$(SSH_TARGET_DIR)
	@echo "### Deployed successfully"

ssh.install: ssh.deploy
	@echo "### Installing debian package on remote machine ..."
	ssh -p 22 $(SSH_USER)@$(SSH_HOST) 'sudo dpkg -i $(SSH_TARGET_DIR)/$(DEB_PACKAGE)'
	@echo "### Installed successfully"



##
## ----------------------
## Clean
## ----------------------
##

clean:
	@echo "### Cleaning ..."
	@rm -f ${BINARY_NAME}
	@rm -f ${BINARY_NAME}_${VERSION}.deb
	@rm -f ./install/${BINARY_NAME}/usr/bin/${BINARY_NAME}
	@rm -f ./install/${BINARY_NAME}/etc/${BINARY_NAME}/config/config.yml
	@echo "### Cleaned successfully"


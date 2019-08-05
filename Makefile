
# General project settings
PROJECT_NAME 		:= efantasy
DOCKER_BASE_IMAGE 	:= $(PROJECT_NAME)-base
SERVICES 		 = $(shell find ./services -name Dockerfile -print0 | xargs -0 -n1 dirname | xargs -n1 basename | sort --unique)
ENVIRONMENT 		?= local

# Versioning
VERSION_LONG 		 = $(shell git describe --first-parent --abbrev=10 --long --tags --dirty)
VERSION_SHORT 		 = $(shell echo $(VERSION_LONG) | cut -f 1 -d "-")
DATE_STRING 		 = $(shell date +'%m-%d-%Y')
GIT_HASH  		 = $(shell git rev-parse --verify HEAD)

# Formatting variables
BOLD 			:= $(shell tput bold)
RESET 			:= $(shell tput sgr0)

.PHONY: services $(SERVICES) docker-base frontend watch tests integration-tests \
	version clean help

.DEFAULT_GOAL := help

services: $(SERVICES)  ## Build Docker image for all services
$(SERVICES):
	@echo "$(BOLD)** Building docker image for service '$@'...$(RESET)"
	docker build -f ./services/$@/Dockerfile -t $(PROJECT_NAME)-$@:latest --build-arg VERSION=$(VERSION_LONG) ./services/$@
	docker tag $(PROJECT_NAME)-$@:latest $(PROJECT_NAME)-$@:$(VERSION_LONG)

docker-base:  ## Build the base image for all services
	@echo "$(BOLD)** Building base image version ${VERSION_LONG}...$(RESET)"
	docker build -f ./Dockerfile -t $(DOCKER_BASE_IMAGE):latest --build-arg VERSION=$(VERSION_LONG) .
	docker tag $(DOCKER_BASE_IMAGE):latest $(DOCKER_BASE_IMAGE):$(VERSION_LONG)

frontend:  ## Build frontend. Requires that the repo is cloned in parent directory to this repo
	@echo "$(BOLD)** Building frontend docker image...$(RESET)"
	docker build -f ../efantasy-frontend/Dockerfile -t $(PROJECT_NAME)-frontend:latest ../efantasy-frontend

watch:  ## Start up a local development environment that watches and redeploys changes automatically
	tilt up --watch

tests:  ## Run all unit tests and print coverage
	go test ./... -v -cover

integration-tests:  ## Run all integration tests, by default against local environment
	python3.6 -m pytest --env $(ENVIRONMENT) -vx -s tests/tests

sec-scan:  ## Run security scan on all repos. Requires 'gosec' installed
	gosec ./...

version:  ## Print the current version
	@echo $(VERSION_LONG)

clean:  ## Clean up all Python cache files and Docker volumes, containers and networks
	@echo "$(BOLD)** Cleaning up Python files...$(RESET)"
	find . -type f -name '*.py[co]' -delete -o -type d -name __pycache__ -delete
	rm -rf .pytest_cache .hypothesis .python-version .mypy_cache
	@echo "$(BOLD)** Cleaning up Docker images and volumes...$(RESET)"
	docker system prune -a -v

help:  ## Print this make target help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@printf "\n"

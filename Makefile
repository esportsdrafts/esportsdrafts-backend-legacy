# General project settings
PROJECT_NAME 		:= esportsdrafts
DOCKER_BASE_IMAGE 	:= $(PROJECT_NAME)-base
SERVICES 		 	 = $(shell find ./services -name Dockerfile -print0 | xargs -0 -n1 dirname | xargs -n1 basename | sort --unique)
ENVIRONMENT 		?= local

# Versioning
VERSION_LONG 		 = $(shell git describe --first-parent --abbrev=10 --long --tags --dirty)
VERSION_SHORT 		 = $(shell echo $(VERSION_LONG) | cut -f 1 -d "-")
DATE_STRING 		 = $(shell date +'%m-%d-%Y')
GIT_HASH  		 	 = $(shell git rev-parse --verify HEAD)

# Formatting variables
BOLD 				:= $(shell tput bold)
RESET 				:= $(shell tput sgr0)

.PHONY: services $(SERVICES) docker-login docker-base docker-test-image frontend \
	watch tests integration-tests version clean help

.DEFAULT_GOAL := help

services: docker-base $(SERVICES)  ## Build Docker image for all services
$(SERVICES):
	@echo "$(BOLD)** Building docker image for service '$@'...$(RESET)"
	docker build -f ./services/$@/Dockerfile -t $(PROJECT_NAME)-$@:latest --build-arg VERSION=$(VERSION_LONG) ./services/$@
	docker tag $(PROJECT_NAME)-$@:latest $(PROJECT_NAME)-$@:$(VERSION_LONG)

docker-login: guard-GH_USERNAME guard-GH_TOKEN  ## Login to github docker registry using $GH_USERNAME and $GH_TOKEN
	@echo "$(BOLD)Logging in to registry $(DOCKER_REGISTRY) ...$(RESET)"
	@docker login docker.pkg.github.com --username ${GH_USERNAME} -p ${GH_TOKEN}

docker-base:  ## Build the base image for all services
	@echo "$(BOLD)** Building base image version ${VERSION_LONG}...$(RESET)"
	docker build -f ./Dockerfile -t $(DOCKER_BASE_IMAGE):latest --build-arg VERSION=$(VERSION_LONG) .
	docker tag $(DOCKER_BASE_IMAGE):latest $(DOCKER_BASE_IMAGE):$(VERSION_LONG)

docker-test-image:  ## Build python-based docker image for integration testing
	@echo "$(BOLD)** Building docker test image version ${VERSION_LONG}...$(RESET)"
	docker build -f ./Dockerfile.testing -t $(PROJECT_NAME)-testing:latest --build-arg VERSION=$(VERSION_LONG) .

frontend:  ## Build frontend. Requires that the repo is cloned in parent directory to this repo
	@echo "$(BOLD)** Building frontend docker image...$(RESET)"
	docker build -f ../esportsdrafts-frontend/Dockerfile -t $(PROJECT_NAME)-frontend:latest ../esportsdrafts-frontend

watch:  ## Start up a local development environment that watches and redeploys changes automatically
	@./scripts/update_minikube_settings.sh
	tilt up --watch

tests:  ## Run all unit tests and print coverage
	go test ./... -v -cover

integration-tests:  ## Run all integration tests, by default against local environment (configure via $ENVIRONMENT)
	python3 -m pytest --env $(ENVIRONMENT) -vx -s tests/tests

sec-scan:  ## Run security scan on all repos. Requires 'gosec' installed
	gosec ./...

version:  ## Print the current version
	@echo $(VERSION_LONG)

clean:  ## Clean up all Python cache files and Docker volumes, containers and networks
	@echo "$(BOLD)** Cleaning up Python files...$(RESET)"
	find . -type f -name '*.py[co]' -delete -o -type d -name __pycache__ -delete
	rm -rf .pytest_cache .hypothesis .python-version .mypy_cache
	@echo "$(BOLD)**Clear out /Users/inbox/...$(RESET)"
	rm -rf /Users/inbox/*
	@echo "$(BOLD)** Cleaning up Docker images and volumes...$(RESET)"
	docker system prune -a

help:  ## Print this make target help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@printf "\n"

guard-%: GUARD
	@ if [ -z '${${*}}' ]; then echo 'Environment variable $(BOLD)$*$(RESET) not set.' && exit 1; fi

# This crap protects against files named the same as the target
.PHONY: GUARD
GUARD:

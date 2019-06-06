
# General project settings
PROJECT_NAME 		:= efantasy
DOCKER_BASE_IMAGE 	:= $(PROJECT_NAME)-base
SERVICES 			 = $(shell find ./services -name Dockerfile -print0 | xargs -0 -n1 dirname | xargs -n1 basename | sort --unique)

# Versioning
VERSION_LONG 		 = $(shell git describe --first-parent --abbrev=10 --long --tags)
VERSION_SHORT 		 = $(shell echo $(VERSION_LONG) | cut -f 1 -d "-")
DATE_STRING 		 = $(shell date +'%m-%d-%Y')
GIT_HASH  			 = $(shell git rev-parse --verify HEAD)

# Formatting variables
BOLD 			:= $(shell tput bold)
RESET 			:= $(shell tput sgr0)

.PHONY: services $(SERVICES) docker-base version help

services: $(SERVICES)  ## Build Docker image for all services
$(SERVICES):
	@echo "$(BOLD)Building docker image for service '$@'...$(RESET)"
	docker build -f ./services/$@/Dockerfile -t --build-arg VERSION=$(VERSION_LONG) $(PROJECT_NAME)-$@:latest ./services/$@
	docker tag $(PROJECT_NAME)-$@:latest $(PROJECT_NAME)-$@:$(VERSION_LONG)

docker-base:  ## Build the base image for all services
	@echo "$(BOLD)** Building base image version ${VERSION_LONG}...$(RESET)"
	docker build -f ./Dockerfile -t $(DOCKER_BASE_IMAGE):latest .
	docker tag $(DOCKER_BASE_IMAGE):latest $(DOCKER_BASE_IMAGE):$(VERSION_LONG)

version:  ## Print the current version
	@echo $(VERSION_LONG)

help:  ## Print this make target help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

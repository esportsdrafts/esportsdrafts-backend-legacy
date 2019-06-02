
# General project settings
PROJECT_NAME 		:= efantasy
DOCKER_BASE_IMAGE 	?= $(PROJECT_NAME)-base
DOCKER_IMAGE 		?= $(PROJECT_NAME)

# Versioning
VERSION_LONG 		 = $(shell git describe --first-parent --abbrev=10 --long --tags)
VERSION_SHORT 		 = $(shell echo $(VERSION_LONG) | cut -f 1 -d "-")
DATE_STRING 		 = $(shell date +'%m-%d-%Y')
GIT_HASH  			 = $(shell git rev-parse --verify HEAD)

# Formatting variables
BOLD 			:= $(shell tput bold)
RESET 			:= $(shell tput sgr0)

.PHONY: version help

version:  ## Print the current version
	@echo $(VERSION_LONG)

help:  ## Print this make target help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
